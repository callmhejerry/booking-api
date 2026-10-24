package dto

type CreateBookingRequest struct {
	SlotID string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type BookingQuery struct {
	Email  string  `form:"email" binding:"required"`
	Limit  uint    `form:"limit,default=10" binding:"min=1,max=50"`
	Cursor *string `form:"cursor,omitempty"`
}
