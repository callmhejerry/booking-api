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
) (db.Slot, error) {
	t.Helper()

	t.Cleanup(func() {
		fmt.Printf("Cleaning up tests\n")
		if err := queries.CleanUpTables(context.Background()); err != nil {
			fmt.Printf("Failed to clean up tables :%s\n", err.Error())
		}
	})

	startTime := time.Now().Add(10 * time.Minute)
	endTime := startTime.Add(30 * time.Minute)

	slot, err := queries.CreateSlot(t.Context(), db.CreateSlotParams{
		SlotID:    "s1",
		Provider:  "Dr Chinedu",
		StartTime: startTime,
		EndTime:   endTime,
	})

	return slot, err
}
