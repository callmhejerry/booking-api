package bookings

import (
	"errors"

	"github.com/callmhejerry/booking-api/internal/slots"
	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
)

func mapSqlErrToAppError(err error) *util.AppError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return slots.ErrSlotNotFound
	}
	if aErr, ok := errors.AsType[*util.AppError](err); ok {
		return aErr
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case util.SQL_UNIQUE_VIOLATION:
			switch pgErr.ConstraintName {
			case "unique_slot_id":
				return slots.ErrSlotAlreadyExist
			case "unique_bookings":
				return slots.ErrSlotAlreadyBooked
			}
		}
	}
	return util.InternalServerError(err)
}
