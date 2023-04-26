package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Accomodation struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Location     string    `json:"location"`
	Wifi         bool      `json:"wifi"`
	Kitchen      bool      `json:"kitchen"`
	AirCondition bool      `json:"airCondition"`
	FreeParking  bool      `json:"freeParking"`
	AutoApproval bool      `json:"autoApproval"`
	Photos       []byte    `json:"photos"`
	MinGuests    uint      `json:"minGuests"`
	MaxGuests    uint      `json:"maxGuests"`
	IDHost       uuid.UUID `json:"id_host"`
}

func (accomodation *Accomodation) BeforeCreate(scope *gorm.DB) error {
	accomodation.ID = uuid.New()
	return nil
}
