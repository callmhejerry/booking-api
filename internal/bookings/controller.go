package bookings

import (
	"net/http"

	"github.com/callmhejerry/booking-api/internal/bookings/dto"
	"github.com/callmhejerry/booking-api/internal/middleware"
	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookingController struct {
	bookingService *BookingService
}

func NewBookingController(bookingService *BookingService, router *gin.Engine) *BookingController {
	c := &BookingController{
		bookingService: bookingService,
	}

	router.Handle("POST", "/bookings", middleware.RateLimitingMiddleware(c.CreateBooking, 2, 5))
	router.Handle("GET", "/bookings", c.ListUserBookings)
	router.Handle("GET", "/bookings/:id", c.GetBookingByID)
	router.Handle("DELETE", "/bookings/:id", c.CancelBooking)

	return c
}

func (c *BookingController) ListUserBookings(ctx *gin.Context) {
	var query dto.BookingQuery
	if err := ctx.BindQuery(&query); err != nil {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid query")
		return
	}

	bookings, err := c.bookingService.ListUserBookings(ctx, query)

	if err != nil {
		util.RespondWithError(ctx, err.HttpStatus, err.Message)
		return
	}

	util.RespondWithSuccess(ctx, bookings)
}

func (c *BookingController) GetBookingByID(ctx *gin.Context) {
	bookingId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, appErr := c.bookingService.GetBookingById(ctx, bookingId)

	if appErr != nil {
		util.RespondWithError(ctx, appErr.HttpStatus, appErr.Message)
		return
	}

	util.RespondWithSuccess(ctx, booking)
}

func (c *BookingController) CreateBooking(ctx *gin.Context) {
	var request dto.CreateBookingRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	idempotencyKey := ctx.GetHeader("Idempotency-Key")

	booking, err := c.bookingService.CreateBooking(ctx, request, idempotencyKey)
	if err != nil {
		util.RespondWithError(ctx, err.HttpStatus, err.Message)
		return
	}

	util.RespondWithCreated(ctx, booking)
}

func (c *BookingController) CancelBooking(ctx *gin.Context) {
	bookingId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, aerr := c.bookingService.CancelBooking(ctx, bookingId)
	if aerr != nil {
		util.RespondWithError(ctx, aerr.HttpStatus, aerr.Message)
		return
	}
	util.RespondWithSuccess(ctx, booking)
}
