package model

type AccomodationSearch struct {
	Location    string `json:"location"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	NumOfGuests int    `json:"numOfGuests"`
}
