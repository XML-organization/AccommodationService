package model

import (
	"time"

	"github.com/google/uuid"
)

type HostGrade struct {
	ID              uuid.UUID `json:"id"`
	AccommodationId uuid.UUID `json:"accommodation_id"`
	UserId          uuid.UUID `json:"user_id"`
	UserName        string    `json:"user_name" `
	UserSurname     string    `json:"user_surname"`
	Grade           float64   `json:"grade"`
	Date            time.Time `json:"date"`
}
