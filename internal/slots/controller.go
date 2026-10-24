package slots

import (
	"net/http"

	"github.com/callmhejerry/booking-api/internal/slots/dto"
	"github.com/callmhejerry/booking-api/internal/util"
	"github.com/gin-gonic/gin"
)

type SlotController struct {
	service *SlotService
}

func NewSlotController(service *SlotService, router *gin.Engine) *SlotController {
	controller := &SlotController{
		service: service,
	}

	router.Handle("POST", "/slots", controller.CreateSlot)
	router.Handle("GET", "/slots", controller.GetProviderSlot)
	router.Handle("POST", "/reset", controller.ResetSlot)
	return controller
}

func (c *SlotController) CreateSlot(ctx *gin.Context) {
	var request dto.CreateSlotRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}

	slot, err := c.service.CreateSlot(ctx, request)
	if err != nil {
		util.RespondWithError(ctx, err.HttpStatus, err.Message)
		return
	}

	util.RespondWithSuccess(ctx, slot)
}

func (c *SlotController) ResetSlot(ctx *gin.Context) {
	var request dto.ResetSlotRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid request body")
		return
	}
	_, err := c.service.ResetSlot(ctx, request)
	if err != nil {
		util.RespondWithError(ctx, err.HttpStatus, err.Message)
		return
	}

	ctx.JSON(200, nil)
}

func (c *SlotController) GetProviderSlot(ctx *gin.Context) {
	provider := ctx.Query("provider")
	if provider == "" {
		util.RespondWithError(ctx, http.StatusBadRequest, "invalid provider")
		return
	}

	slots, err := c.service.GetSlotForProvider(ctx, provider)
	if err != nil {
		util.RespondWithError(ctx, err.HttpStatus, err.Message)
		return
	}

	util.RespondWithSuccess(ctx, slots)
}
