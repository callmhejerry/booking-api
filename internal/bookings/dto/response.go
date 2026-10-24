package dto

type PaginatedBookingResponse struct {
	Bookings   []BookingResponseDto `json:"bookings"`
	NextCursor *string              `json:"next_cursor"`
}

type BookingResponseDto struct {
	BookingID string `json:"id"`
	SlotID    string `json:"slot_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
