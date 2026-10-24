package slots

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/callmhejerry/booking-api/internal/slots/dto"
	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/gin-gonic/gin"
)

type SlotService struct {
	repo   SlotRepository
	logger *slog.Logger
}

func NewSlotService(repo SlotRepository, logger *slog.Logger) *SlotService {
	return &SlotService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SlotService) CreateSlot(ctx context.Context, request dto.CreateSlotRequest) (*dto.SlotResponseDto, *util.AppError) {
	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		s.logger.Error("invalid start time", slog.String("start_time", request.StartTime))
		return nil, &util.AppError{
			Message:    "invalid start_time",
			HttpStatus: http.StatusBadRequest,
		}
	}
	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		s.logger.Error("invalid start time", slog.String("end_time", request.EndTime))
		return nil, &util.AppError{
			Message:    "invalid end_time",
			HttpStatus: http.StatusBadRequest,
		}
	}

	if startTime.After(endTime) {
		s.logger.Error("Start time is after end time",
			slog.String("start_time", request.StartTime),
			slog.String("end_time", request.EndTime),
		)
		return nil, &util.AppError{
			Message:    "Start time cannot be after end time",
			HttpStatus: http.StatusBadRequest,
		}
	}

	slot, appErr := s.repo.CreateSlot(
		ctx, request.SlotID, request.Provider, startTime, endTime,
	)
	if appErr != nil {
		s.logger.Error(appErr.Error(),
			slog.String("error", appErr.Err.Error()),
			slog.String("start_time", request.StartTime),
			slog.String("end_time", request.EndTime),
			slog.String("provider", request.Provider),
			slog.String("slot_id", request.SlotID),
		)
		return nil, appErr
	}
	return &dto.SlotResponseDto{
		ID:        slot.SlotID,
		Provider:  slot.Provider,
		Start:     slot.StartTime.Format(time.RFC3339),
		End:       slot.EndTime.Format(time.RFC3339),
		Available: slot.IsAvailable,
	}, nil
}

func (s *SlotService) ResetSlot(
	ctx context.Context,
	request dto.ResetSlotRequest,
) ([]dto.SlotResponseDto, *util.AppError) {
	var errs []*util.AppError
	slots := make([]dto.SlotResponseDto, 0, len(request.Slots))

	if err := s.repo.DeleteAllSlots(ctx); err != nil {
		s.logger.Error("Failed to delete all slots",
			"error", err.Error(),
		)
		return nil, &util.AppError{
			Message:    "Failed to reset all slots",
			HttpStatus: http.StatusInternalServerError,
		}
	}

	for i := range request.Slots {
		slot, err := s.CreateSlot(ctx, request.Slots[i])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		slots = append(slots, *slot)
	}

	if len(slots) == 0 {
		return nil, errs[0]
	}
	return slots, nil
}

func (s *SlotService) GetSlotForProvider(
	ctx *gin.Context,
	provider string,
) (*dto.ListSlotResponseDto, *util.AppError) {
	slots, err := s.repo.GetSlotForProvider(ctx, provider)

	if err != nil {
		s.logger.Error("Failed to get slots for provider",
			"error", err.Err.Error(),
			"provider", provider,
		)
		return nil, err
	}

	response := make([]dto.SlotResponseDto, len(slots))

	for i := range slots {
		response[i] = dto.SlotResponseDto{
			ID:        slots[i].SlotID,
			Provider:  slots[i].Provider,
			Start:     slots[i].StartTime.Format(time.RFC3339),
			End:       slots[i].EndTime.Format(time.RFC3339),
			Available: slots[i].IsAvailable,
		}
	}

	return &dto.ListSlotResponseDto{
		Slot: response,
	}, nil
}
