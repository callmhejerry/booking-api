package bookings

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/callmhejerry/booking-api/internal/bookings/dto"
	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/slots"
	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingService struct {
	repo     BookingRepository
	slotRepo slots.SlotRepository
	cache    *BookingCache
	logger   *slog.Logger
}

func NewBookingService(
	repo BookingRepository,
	slotRepo slots.SlotRepository,
	bookingCache *BookingCache,
	logger *slog.Logger,
) *BookingService {

	return &BookingService{
		repo:     repo,
		logger:   logger,
		slotRepo: slotRepo,
		cache:    bookingCache,
	}
}

func (s *BookingService) CreateBooking(
	ctx *gin.Context,
	request dto.CreateBookingRequest,
	idempotencyKey string,
) (*dto.BookingResponseDto, *util.AppError) {
	requestHash, err := util.HashStruct(request)
	if err != nil {
		return nil, err
	}

	if idempotencyKey != "" {
		cachedData, err := s.cache.Get(idempotencyKey, requestHash)
		if err != nil {
			s.logger.Error(err.Message)
			return nil, err
		}
		if cachedData != nil {
			return cachedData, nil
		}
	}

	slot, err := s.slotRepo.GetSlotBySlotId(ctx, request.SlotID)

	if err != nil {
		return nil, err
	}

	if slot.StartTime.Before(time.Now()) {
		errMsg := "Slot has already started, please book another slot"
		s.logger.Error(errMsg,
			"start_time", slot.StartTime.Format(time.RFC3339),
			"end_time", slot.EndTime.Format(time.RFC3339),
		)
		return nil, &util.AppError{
			Message:    errMsg,
			HttpStatus: http.StatusBadRequest,
			Err:        nil,
		}
	}

	booked, err := s.repo.CreateBooking(
		ctx,
		request.SlotID,
		request.Name,
		request.Email,
	)
	if err != nil {
		s.logger.Error("Failed to create booking",
			"error", err.Error(),
			"slot_id", request.SlotID,
			"name", request.Name,
			"email", request.Email,
		)
		return nil, &util.AppError{
			Message:    "Failed to create booking",
			HttpStatus: http.StatusBadRequest,
		}
	}

	response := convertToBookingResponse(*booked)
	if idempotencyKey != "" {
		s.cache.Add(idempotencyKey, requestHash, response)
	}
	return &response, nil
}

func (s *BookingService) ListUserBookings(ctx *gin.Context, query dto.BookingQuery) (*dto.PaginatedBookingResponse, *util.AppError) {
	slotId, createdAt, err := convertFromBookingCursor(query.Cursor)

	if err != nil {
		errMsg := "Failed to convert from booking cursor"
		s.logger.Error(errMsg, "error", err.Error(), "cursor", query.Cursor)
		return nil, &util.AppError{
			Message:    errMsg,
			HttpStatus: http.StatusBadRequest,
		}
	}
	bookings, err := s.repo.ListUserBookings(
		ctx, query.Email,
		slotId, createdAt, query.Limit+1,
	)

	if err != nil {
		errMsg := "internal server error"
		s.logger.Error(errMsg, "error", err.Error())
		return nil, &util.AppError{
			Message:    errMsg,
			HttpStatus: http.StatusInternalServerError,
		}
	}

	// There are no bookings left
	if len(bookings) <= int(query.Limit) {
		return &dto.PaginatedBookingResponse{
			Bookings:   convertToManyBookingResponse(bookings),
			NextCursor: nil,
		}, nil
	}

	lastBooking := bookings[len(bookings)-1]
	lastBookingSlotId := lastBooking.SlotID
	lastBookingCreatedAt := lastBooking.CreatedAt.Format(time.RFC3339)

	nextCursor := convertToBookingCursor(lastBookingSlotId, lastBookingCreatedAt)

	return &dto.PaginatedBookingResponse{
		Bookings:   convertToManyBookingResponse(bookings[:len(bookings)-1]),
		NextCursor: &nextCursor,
	}, nil
}

func (s *BookingService) GetBookingById(ctx *gin.Context, id uuid.UUID) (*dto.BookingResponseDto, *util.AppError) {
	booking, err := s.repo.GetBookingById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get booking", "error", err.Error(), "booking_id", id)
		return nil, &util.AppError{
			Message:    "slot not found",
			HttpStatus: http.StatusNotFound,
		}
	}
	return &dto.BookingResponseDto{
		BookingID: booking.ID.String(),
		SlotID:    booking.SlotID,
		Name:      booking.Name,
		Email:     booking.Email,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *BookingService) CancelBooking(ctx *gin.Context, id uuid.UUID) (*dto.BookingResponseDto, *util.AppError) {
	cancelledBooking, err := s.repo.CancelBooking(ctx, id)
	if err != nil {
		slog.Error("Failed to cancel booking",
			"error", err.Error(),
			"booking_id", id,
		)
		return nil, &util.AppError{
			Message:    "Failed to cancel booking",
			HttpStatus: http.StatusInternalServerError,
		}
	}

	response := convertToBookingResponse(*cancelledBooking)
	return &response, nil
}

func convertToBookingResponse(booking db.Booking) dto.BookingResponseDto {
	return dto.BookingResponseDto{
		BookingID: booking.ID.String(),
		SlotID:    booking.SlotID,
		Name:      booking.Name,
		Email:     booking.Email,
		Status:    booking.Status,
		CreatedAt: booking.CreatedAt.Format(time.RFC3339),
	}
}

func convertToManyBookingResponse(bookings []db.Booking) []dto.BookingResponseDto {
	bookingResponse := make([]dto.BookingResponseDto, 0, len(bookings))

	for _, booking := range bookings {
		bookingResponse = append(bookingResponse, dto.BookingResponseDto{
			BookingID: booking.ID.String(),
			SlotID:    booking.SlotID,
			Name:      booking.Name,
			Email:     booking.Email,
			Status:    booking.Status,
			CreatedAt: booking.CreatedAt.Format(time.RFC3339),
		})
	}

	return bookingResponse
}

func convertFromBookingCursor(cursor *string) (*string, *string, error) {
	if cursor == nil {
		return nil, nil, nil
	}
	decodedCursorBytes, err := base64.URLEncoding.DecodeString(*cursor)
	if err != nil {
		return nil, nil, err
	}
	decodedCursor := string(decodedCursorBytes)

	var slotId string
	var createdAt string

	_, err = fmt.Sscanf(decodedCursor, "%s:%s", &slotId, &createdAt)
	if err != nil {
		return nil, nil, err
	}
	return &slotId, &createdAt, nil
}

func convertToBookingCursor(slotId, createdAt string) string {
	cursorString := fmt.Sprintf("%s:%s", slotId, createdAt)
	return base64.URLEncoding.EncodeToString([]byte(cursorString))
}
