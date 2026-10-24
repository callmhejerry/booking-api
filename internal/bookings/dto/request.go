package dto

type CreateBookingRequest struct {
	SlotID string `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

type BookingQuery struct {
	Email  string  `form:"email" binding:"required"`
	Limit  uint    `form:"limit,default=10" binding:"min=1,max=50"`
	Cursor *string `form:"cursor,omitempty"`
}
