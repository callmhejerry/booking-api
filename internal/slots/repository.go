package slots

import (
	"context"
	"time"

	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/util"
)

type SlotRepository interface {
	CreateSlot(ctx context.Context, slotId string, provider string,
		startTime, endTime time.Time) (*db.Slot, *util.AppError)
	ListSlots(ctx context.Context) ([]db.Slot, error)
	DeleteAllSlots(ctx context.Context) error
	GetSlotForProvider(ctx context.Context, provider string) ([]db.Slot, *util.AppError)
	GetSlotBySlotId(ctx context.Context, slotId string) (*db.Slot, *util.AppError)
}

type SlotRepositoryImpl struct {
	db *db.Queries
}

func NewRepositoryImpl(db *db.Queries) *SlotRepositoryImpl {
	return &SlotRepositoryImpl{
		db: db,
	}
}

func (r *SlotRepositoryImpl) CreateSlot(
	ctx context.Context,
	slotID string,
	provider string,
	startTime, endTime time.Time,
) (*db.Slot, *util.AppError) {
	slot, err := r.db.CreateSlot(ctx, db.CreateSlotParams{
		Provider:  provider,
		StartTime: startTime,
		EndTime:   endTime,
	})
	return &slot, mapSqlErrToAppError(err)
}

func (r *SlotRepositoryImpl) ListSlots(ctx context.Context) ([]db.Slot, error) {
	return r.db.ListSlots(ctx)
}

func (r *SlotRepositoryImpl) DeleteAllSlots(ctx context.Context) error {
	return r.db.ResetSlots(ctx)
}

func (r *SlotRepositoryImpl) GetSlotForProvider(ctx context.Context, provider string) ([]db.Slot, *util.AppError) {
	slot, err := r.db.GetProviderSlot(ctx, provider)
	if err != nil {
		return nil, mapSqlErrToAppError(err)
	}

	return slot, nil
}

func (r *SlotRepositoryImpl) GetSlotBySlotId(
	ctx context.Context,
	slotId string,
) (*db.Slot, *util.AppError) {
	slot, err := r.db.GetSlotBySlotId(ctx, slotId)
	if err != nil {
		return nil, mapSqlErrToAppError(err)
	}
	return &slot, nil
}
