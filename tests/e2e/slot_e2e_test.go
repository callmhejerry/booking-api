package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/callmhejerry/booking-api/internal/slots/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateSlot(t *testing.T) {
	payload := dto.CreateSlotRequest{
		SlotID:    "s1",
		Provider:  "Dr Chinedu",
		StartTime: time.Now().Format(time.RFC3339),
		EndTime:   time.Now().Add(time.Minute * 30).Format(time.RFC3339),
	}

	jsonBody, _ := json.Marshal(payload)

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/slots", App.TestServerUrl), bytes.NewBuffer(jsonBody))
	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)
	defer request.Body.Close()

	response, err := http.DefaultClient.Do(request)

	require.NoError(t, err)

	defer response.Body.Close()

	assert.Equal(t, response.StatusCode, http.StatusOK)
}

func Test_ResetSlot(t *testing.T) {
	slots := [3]dto.CreateSlotRequest{
		{
			SlotID:    "s2",
			Provider:  "Dr s2",
			StartTime: time.Now().Add(3 * time.Minute).Format(time.RFC3339),
			EndTime:   time.Now().Add((3 + 10) * time.Minute).Format(time.RFC3339),
		},
		{
			SlotID:    "s3",
			Provider:  "Dr s3",
			StartTime: time.Now().Add(13 * time.Minute).Format(time.RFC3339),
			EndTime:   time.Now().Add((13 + 10) * time.Minute).Format(time.RFC3339),
		},
		{
			SlotID:    "s4",
			Provider:  "Dr s4",
			StartTime: time.Now().Add(23 * time.Minute).Format(time.RFC3339),
			EndTime:   time.Now().Add((13 + 10) * time.Minute).Format(time.RFC3339),
		},
	}

	payload := dto.ResetSlotRequest{
		Slots: slots[:],
	}
	jsonBytes, _ := json.Marshal(payload)

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/reset", App.TestServerUrl), bytes.NewBuffer(jsonBytes))
	defer request.Body.Close()

	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func Test_GetSlotById(t *testing.T) {
	endpoint, _ := url.Parse(fmt.Sprintf("%s/slots", App.TestServerUrl))
	query := url.Values{}
	query.Set("provider", "Dr Chinedu")

	endpoint.RawQuery = query.Encode()

	request, _ := http.NewRequest("GET", endpoint.String(), nil)
	request.Header.Set("X-API-Key", App.Config.AppCongfig.ApiKey)
	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	defer response.Body.Close()

	var actualResponse dto.ListSlotResponseDto

	_ = json.NewDecoder(response.Body).Decode(&actualResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	// assert.Equal(t, 1, len(actualResponse.Slot))
	// assert.True(t, actualResponse.Slot[0].Available)
}
