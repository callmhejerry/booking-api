package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/callmhejerry/booking-api/internal/db"
)

func SeedSlot(
	t *testing.T,
	queries *db.Queries,
	startTime time.Time,
	endTime time.Time,
) (db.Slot, error) {
	t.Helper()

	t.Cleanup(func() {
		fmt.Printf("Cleaning up tests\n")
		if err := queries.CleanUpTables(context.Background()); err != nil {
			fmt.Printf("Failed to clean up tables :%s\n", err.Error())
		}
	})

	slot, err := queries.CreateSlot(t.Context(), db.CreateSlotParams{
		SlotID:    "s1",
		Provider:  "Dr Chinedu",
		StartTime: startTime,
		EndTime:   endTime,
	})

	return slot, err
}
