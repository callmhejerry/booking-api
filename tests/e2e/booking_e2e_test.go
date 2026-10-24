package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/callmhejerry/booking-api/internal/bookings/dto"
	"github.com/callmhejerry/booking-api/tests/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateBooking(t *testing.T) {
	url := fmt.Sprintf("%s/bookings", App.TestServerUrl)

	slot, err := helpers.SeedSlot(t, App.Queries, time.Now().Add(5*time.Minute), time.Now().Add(15*time.Minute))
	require.NoError(t, err)

	payload := dto.CreateBookingRequest{
		SlotID: slot.SlotID,
		Name:   "Kelechi",
		Email:  "kelechichineduben@gmail.com",
	}
	jsonBody, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)

	res, err := http.DefaultClient.Do(request)
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
	slot, err := helpers.SeedSlot(t, App.Queries, time.Now().Add(5*time.Minute), time.Now().Add(time.Minute*15))
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

func Test_CreateBookingTimePassed(t *testing.T) {
	slot, err := helpers.SeedSlot(t, App.Queries, time.Now(), time.Now().Add(time.Minute*15))
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

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// func Test_CreateBookingRateLimit(t *testing.T) {
// 	slot, err := helpers.SeedSlot(t, App.Queries, time.Now().Add(5*time.Minute), time.Now().Add(time.Minute*15))
// 	require.NoError(t, err)

// 	payload := dto.CreateBookingRequest{
// 		SlotID: slot.SlotID,
// 		Name:   "Kelechi",
// 		Email:  "kelechichineduben@gmail.com",
// 	}

// 	header := make(map[string]string)
// 	idempotencyKey := uuid.New()
// 	header["Idempotency-Key"] = idempotencyKey.String()

// 	fmt.Printf("Idempotency key is %s\n", idempotencyKey.String())

// 	for range 5 {
// 		createBooking(t, payload, header)
// 		time.Sleep(10 * time.Millisecond)
// 	}

// 	res := createBooking(t, payload, header)

// 	assert.Equal(t, http.StatusTooManyRequests, res.StatusCode)
// }

func Test_GetProviderBookings(t *testing.T) {
	url, err := url.Parse(fmt.Sprintf("%s/bookings", App.TestServerUrl))
	require.NoError(t, err)
	query := url.Query()
	query.Set("email", "kelechichineduben@gmail.com")
	url.RawQuery = query.Encode()

	request, err := http.NewRequest("GET", url.String(), nil)
	require.NoError(t, err)

	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	var actualResponse dto.PaginatedBookingResponse
	err = json.NewDecoder(response.Body).Decode(&actualResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func Test_CancelBooking(t *testing.T) {
	url := fmt.Sprintf("%s/bookings", App.TestServerUrl)

	slot, err := helpers.SeedSlot(
		t,
		App.Queries,
		time.Now().Add(5*time.Minute),
		time.Now().Add(15*time.Minute),
	)
	require.NoError(t, err)

	payload := dto.CreateBookingRequest{
		SlotID: slot.SlotID,
		Name:   "Kelechi",
		Email:  "kelechichineduben@gmail.com",
	}

	jsonBody, err := json.Marshal(payload)
	require.NoError(t, err)

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(jsonBody),
	)
	require.NoError(t, err)

	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)
	request.Header.Set("Content-Type", "application/json")

	time.Sleep(time.Second * 1)
	res, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	require.NotNil(t, res)
	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)

	var actualResponse dto.BookingResponseDto
	err = json.NewDecoder(res.Body).Decode(&actualResponse)
	require.NoError(t, err)
	require.NotEmpty(t, actualResponse.BookingID)

	deleteURL := fmt.Sprintf("%s/%s", url, actualResponse.BookingID)

	cancelRequest, err := http.NewRequest(
		http.MethodDelete,
		deleteURL,
		nil,
	)
	require.NoError(t, err)

	cancelRequest.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)

	cancelResponse, err := http.DefaultClient.Do(cancelRequest)
	require.NoError(t, err)
	require.NotNil(t, cancelResponse)
	defer cancelResponse.Body.Close()

	require.Equal(t, http.StatusOK, cancelResponse.StatusCode)

	var deleteResponse dto.BookingResponseDto
	err = json.NewDecoder(cancelResponse.Body).Decode(&deleteResponse)
	require.NoError(t, err)

	assert.Equal(t, "cancelled", deleteResponse.Status)
}

func Test_GetBookingByID(t *testing.T) {
	url := fmt.Sprintf("%s/bookings", App.TestServerUrl)

	slot, err := helpers.SeedSlot(
		t,
		App.Queries,
		time.Now().Add(5*time.Minute),
		time.Now().Add(15*time.Minute),
	)
	require.NoError(t, err)

	payload := dto.CreateBookingRequest{
		SlotID: slot.SlotID,
		Name:   "Kelechi",
		Email:  "kelechichineduben@gmail.com",
	}

	jsonBody, err := json.Marshal(payload)
	require.NoError(t, err)

	request, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(jsonBody),
	)
	require.NoError(t, err)

	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)
	request.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	require.NotNil(t, res)
	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)

	var actualResponse dto.BookingResponseDto
	err = json.NewDecoder(res.Body).Decode(&actualResponse)
	require.NoError(t, err)
	require.NotEmpty(t, actualResponse.BookingID)

	getURl := fmt.Sprintf("%s/%s", url, actualResponse.BookingID)

	getRequest, err := http.NewRequest(
		http.MethodGet,
		getURl,
		nil,
	)
	require.NoError(t, err)

	getRequest.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)

	getResponse, err := http.DefaultClient.Do(getRequest)
	require.NoError(t, err)
	require.NotNil(t, getResponse)
	defer getResponse.Body.Close()

	require.Equal(t, http.StatusOK, getResponse.StatusCode)

	var deleteResponse dto.BookingResponseDto
	err = json.NewDecoder(getResponse.Body).Decode(&deleteResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, getResponse.StatusCode)
}

func createBooking(
	t *testing.T,
	payload dto.CreateBookingRequest,
	headers map[string]string,
) *http.Response {
	t.Helper()
	url := fmt.Sprintf("%s/bookings", App.TestServerUrl)
	jsonBody, _ := json.Marshal(payload)

	client := http.DefaultClient

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	res, err := client.Do(request)
	require.NoError(t, err)
	return res
}
