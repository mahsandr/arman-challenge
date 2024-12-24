package models

type UserSegment struct {
	UserID      string `json:"user_id" validate:"required"`
	Segment     string `json:"segment" validate:"required"`
	RegistredAt uint32
}
