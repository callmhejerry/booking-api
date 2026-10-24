package slots

import (
	"errors"
	"net/http"

	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/jackc/pgx"
)

var (
	ErrSlotAlreadyExist  *util.AppError = util.New("Slot with this slot_id already exist", http.StatusConflict, nil)
	ErrSlotNotFound      *util.AppError = util.New("slot not found", http.StatusNotFound, nil)
	ErrSlotAlreadyBooked *util.AppError = util.New("Slot has already been booked", http.StatusConflict, nil)
)

func mapSqlErrToAppError(err error) *util.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrSlotNotFound
	}

	if pgErr, ok := errors.AsType[pgx.PgError](err); ok {
		switch pgErr.Code {
		case util.SQL_UNIQUE_VIOLATION:
			return ErrSlotAlreadyExist
		}
	}
	return util.InternalServerError(err)
}
