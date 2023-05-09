package handler

import (
	"accomodation-service/model"

	pb "github.com/XML-organization/common/proto/accomodation_service"
	"github.com/google/uuid"
)

func mapAccomodationFromCreateAccomodation(accomodation *pb.CreateRequest) model.Accomodation {
	accomodationID, err := uuid.Parse(accomodation.ID)
	if err != nil {
		panic(err)
	}

	iDHost, err := uuid.Parse(accomodation.IDHost)
	if err != nil {
		panic(err)
	}
	return model.Accomodation{
		ID:           accomodationID,
		Name:         accomodation.Name,
		Location:     accomodation.Location,
		Wifi:         accomodation.Wifi,
		Kitchen:      accomodation.Kitchen,
		AirCondition: accomodation.AirCondition,
		FreeParking:  accomodation.FreeParking,
		AutoApproval: accomodation.AutoApproval,
		Photos:       accomodation.Photos,
		MinGuests:    uint(accomodation.MinGuests),
		MaxGuests:    uint(accomodation.MaxGuests),
		IDHost:       iDHost,
	}
}
