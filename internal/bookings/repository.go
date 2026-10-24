package bookings

import (
	"errors"

	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/slots"
	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingStatus string

const (
	BOOKING_CANCELLED_STATUS BookingStatus = "cancelled"
	BOOKING_CONFIRMED_STATUS BookingStatus = "confirmed"
)

type BookingRepository interface {
	CreateBooking(
		ctx *gin.Context,
		slotId, name, email string,
	) (*db.Booking, *util.AppError)

	ListUserBookings(
		ctx *gin.Context,
		email string, slotId, createdAt *string,
		limit uint,
	) ([]db.Booking, error)

	GetBookingById(ctx *gin.Context, id uuid.UUID) (*db.Booking, error)
	CancelBooking(ctx *gin.Context, id uuid.UUID) (*db.Booking, error)
}

type BookingRepositoryImpl struct {
	db     *db.Queries
	dbPool *pgxpool.Pool
}

func NewBookingRepositoryImpl(db *db.Queries, dbPool *pgxpool.Pool) *BookingRepositoryImpl {
	return &BookingRepositoryImpl{
		db:     db,
		dbPool: dbPool,
	}
}

func (r *BookingRepositoryImpl) CreateBooking(
	ctx *gin.Context,
	slotId, name, email string,
) (*db.Booking, *util.AppError) {
	var bookingResult *db.Booking

	err := db.WithTx(ctx, r.dbPool, func(qtx *db.Queries) error {
		bookedSlot, err := r.db.BookSlot(ctx, slotId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return slots.ErrSlotAlreadyBooked
			}
			return err
		}
		bookings, err := r.db.CreateBooking(ctx, db.CreateBookingParams{
			SlotID: bookedSlot.SlotID,
			Name:   name,
			Email:  email,
			Status: string(BOOKING_CONFIRMED_STATUS),
		})
		if err != nil {
			return err
		}

		bookingResult = &bookings
		return nil
	})

	if err != nil {
		return nil, mapSqlErrToAppError(err)
	}

	return bookingResult, nil
}

func (r *BookingRepositoryImpl) ListUserBookings(
	ctx *gin.Context,
	email string,
	slotId, createdAt *string,
	limit uint,
) ([]db.Booking, error) {
	userBookings, err := r.db.GetUserBooking(ctx, db.GetUserBookingParams{
		Email:     email,
		SlotID:    slotId,
		CreatedAt: createdAt,
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	return userBookings, nil
}

func (r *BookingRepositoryImpl) GetBookingById(ctx *gin.Context, id uuid.UUID) (*db.Booking, error) {
	booking, err := r.db.GetBookingById(ctx, id)

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) CancelBooking(
	ctx *gin.Context,
	id uuid.UUID,
) (*db.Booking, error) {
	var bookingResult *db.Booking
	err := db.WithTx(ctx, r.dbPool, func(qtx *db.Queries) error {
		booking, err := qtx.CancelBooking(ctx, db.CancelBookingParams{
			ID:     id,
			Status: string(BOOKING_CANCELLED_STATUS),
		})
		if err != nil {
			return err
		}
		_, err = qtx.UnBookSlot(ctx, booking.SlotID)

		if err != nil {
			return err
		}

		bookingResult = &booking

		return nil
	})

	if err != nil {
		return nil, err
	}

	return bookingResult, nil
}
