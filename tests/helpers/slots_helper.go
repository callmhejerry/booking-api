package helpers

import (
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
		fmt.Printf("Cleaning up tests")
		queries.CleanUpTables(t.Context())
	})

	startTime := time.Now().Add(10 * time.Minute)
	endTime := startTime.Add(30 * time.Minute)

	return queries.CreateSlot(t.Context(), db.CreateSlotParams{
		SlotID:    "s1",
		Provider:  "Dr Chinedu",
		StartTime: startTime,
		EndTime:   endTime,
	})
}
