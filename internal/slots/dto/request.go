package dto

type CreateSlotRequest struct {
	SlotID    string `json:"id"`
	Provider  string `json:"provider"`
	StartTime string `json:"start"`
	EndTime   string `json:"end"`
}

type ResetSlotRequest struct {
	Slots []CreateSlotRequest `json:"slots"`
}
