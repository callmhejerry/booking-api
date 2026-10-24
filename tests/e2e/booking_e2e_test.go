package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/callmhejerry/booking-api/internal/bookings/dto"
	"github.com/callmhejerry/booking-api/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateBooking(t *testing.T) {
	url := fmt.Sprintf("%s/bookings", TestServerURL)

	slot, err := helpers.SeedSlot(t, &Queries)
	require.NoError(t, err)

	payload := dto.CreateBookingRequest{
		SlotID: slot.SlotID,
		Name:   "Kelechi",
		Email:  "kelechichineduben@gmail.com",
	}
	jsonBody, _ := json.Marshal(payload)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)

	defer res.Body.Close()

	var actualResponse dto.BookingResponseDto
	err = json.NewDecoder(res.Body).Decode(&actualResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, payload.SlotID, actualResponse.SlotID)
	assert.Equal(t, "confirmed", actualResponse.Status)
}

func Test_CreateBookingWithIdempotencySuccess(
	t *testing.T,
) {

	slot, err := helpers.SeedSlot(t, &Queries)
	require.NoError(t, err)

	payload := dto.CreateBookingRequest{
		SlotID: slot.SlotID,
		Name:   "Kelechi",
		Email:  "kelechichineduben@gmail.com",
	}

	header := make(map[string]string)
	idempotencyKey := uuid.New()
	header["Idempotency-Key"] = idempotencyKey.String()

	fmt.Printf("Idempotency key is %s\n", idempotencyKey.String())

	res := createBooking(t, payload, header)
	defer res.Body.Close()

	var actualResponse dto.BookingResponseDto
	_ = json.NewDecoder(res.Body).Decode(&actualResponse)

	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, payload.SlotID, actualResponse.SlotID)
	assert.Equal(t, "confirmed", actualResponse.Status)

	res2 := createBooking(t, payload, header)
	var actualResponse2 dto.BookingResponseDto
	_ = json.NewDecoder(res2.Body).Decode(&actualResponse2)

	assert.Equal(t, actualResponse, actualResponse2)
}

func createBooking(
	t *testing.T,
	payload dto.CreateBookingRequest,
	headers map[string]string,
) *http.Response {
	t.Helper()
	url := fmt.Sprintf("%s/bookings", TestServerURL)
	jsonBody, _ := json.Marshal(payload)

	client := http.DefaultClient

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	res, err := client.Do(request)
	require.NoError(t, err)
	return res
}
