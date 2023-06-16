package model

import (
	"github.com/google/uuid"
)

type RecommendedAccommodation struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Location      string    `json:"location"`
	Wifi          bool      `json:"wifi"`
	Kitchen       bool      `json:"kitchen"`
	AirCondition  bool      `json:"airCondition"`
	FreeParking   bool      `json:"freeParking"`
	AutoApproval  bool      `json:"autoApproval"`
	PricePerGuest bool      `json:"pricePerGuest"`
	Photos        []byte    `json:"photos"`
	MinGuests     int       `json:"minGuests"`
	MaxGuests     int       `json:"maxGuests"`
	IDHost        uuid.UUID `gorm:"column:id_host" json:"idHost"`
	AverageGrade  float64   `json:"averageGrade"`
}
