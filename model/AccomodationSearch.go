package model

import "time"

type AccomodationSearch struct {
	Location    string    `json:"location"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	NumOfGuests int       `json:"numOfGuests"`
}
