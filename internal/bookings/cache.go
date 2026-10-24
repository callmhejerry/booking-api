package bookings

import (
	"sync"

	"github.com/callmhejerry/booking-api/internal/bookings/dto"
	"github.com/callmhejerry/booking-api/internal/util"
)

type BookingCache struct {
	cache map[string]BookingCacheData
	mu    sync.Mutex
}

type BookingCacheData struct {
	response    dto.BookingResponseDto
	requestHash string
}

func NewBookingCache() *BookingCache {
	return &BookingCache{
		cache: make(map[string]BookingCacheData),
	}
}

func (c *BookingCache) Add(
	key string,
	requestHash string,
	data dto.BookingResponseDto,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheData := BookingCacheData{
		response:    data,
		requestHash: requestHash,
	}

	c.cache[key] = cacheData
}

func (c *BookingCache) Get(
	key string,
	requestHash string,
) (*dto.BookingResponseDto, *util.AppError) {
	if cachedData, ok := c.cache[key]; ok {
		// we have the data in cache
		if requestHash != cachedData.requestHash {
			return nil, &util.AppError{
				HttpStatus: 422,
				Message:    "Same idempotency key was used with a different body",
				Err:        nil,
			}
		}
		return &cachedData.response, nil
	}
	// we dont have the data in cache
	return nil, nil
}
