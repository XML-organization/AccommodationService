package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Availability struct {
	ID             uuid.UUID `json:"id"`
	StartDate      string    `json:"start_date"`
	EndDate        string    `json:"end_date" `
	IdAccomodation uuid.UUID `json:"id_accomodation"`
	Price          float64   `json:"price"`
}

func (availability *Availability) BeforeCreate(scope *gorm.DB) error {
	availability.ID = uuid.New()
	return nil
}
