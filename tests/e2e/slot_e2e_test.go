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

	response, err := http.Post(fmt.Sprintf("%s/slots", TestServerURL), "application/json", bytes.NewBuffer(jsonBody))

	require.NoError(t, err)

	defer response.Body.Close()

	assert.Equal(t, response.StatusCode, http.StatusOK)
}

func Test_GetSlotById(t *testing.T) {

	endpoint, _ := url.Parse(fmt.Sprintf("%s/slots", TestServerURL))
	query := url.Values{}
	query.Set("provider", "Dr Chinedu")

	endpoint.RawQuery = query.Encode()

	response, err := http.Get(endpoint.String())
	require.NoError(t, err)

	defer response.Body.Close()

	var actualResponse dto.ListSlotResponseDto

	_ = json.NewDecoder(response.Body).Decode(&actualResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	// assert.Equal(t, 1, len(actualResponse.Slot))
	// assert.True(t, actualResponse.Slot[0].Available)
}
