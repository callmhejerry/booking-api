package dto

type SlotResponseDto struct {
	ID        string `json:"id"`
	Provider  string `json:"provider"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Available bool   `json:"available"`
}

type ListSlotResponseDto struct {
	Slot []SlotResponseDto `json:"slots"`
}
