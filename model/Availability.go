package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Availability struct {
	ID             uuid.UUID `json:"id"`
	StartDate      string    `json:"startDate" gorm:"not null;type:string"`
	EndDate        string    `json:"endDate" gorm:"not null;type:string"`
	IdAccomodation uuid.UUID `json:"id_accomodation"`
}

func (availability *Availability) BeforeCreate(scope *gorm.DB) error {
	availability.ID = uuid.New()
	return nil
}
