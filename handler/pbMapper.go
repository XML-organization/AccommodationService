package handler

import (
	"accomodation-service/model"
	"fmt"
	"strconv"

	pb "github.com/XML-organization/common/proto/accomodation_service"
	"github.com/google/uuid"
)

func mapAccomodationFromCreateAccomodation(accomodation *pb.CreateRequest) model.Accomodation {
	/* 	accomodationID, err := uuid.Parse(accomodation.ID)
	   	if err != nil {
	   		panic(err)
	   	}

	iDHost, err := uuid.Parse(accomodation.IDHost)
	if err != nil {
		panic(err)
	}*/
	println(accomodation.IDHost)
	hostID, err := uuid.Parse(accomodation.IDHost)
	if err != nil {
		panic(err)
	}
	min, err := strconv.Atoi(accomodation.MinGuests)
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(accomodation.MaxGuests)
	if err != nil {
		panic(err)
	}
	/* 	wifi, err := strconv.ParseBool(accomodation.Wifi)
	   	if err != nil {
	   		fmt.Println("Error parsing bool:", err)
	   		panic(err)
	   	}
	   	kitchen, err := strconv.ParseBool(accomodation.Kitchen)
	   	if err != nil {
	   		fmt.Println("Error parsing bool:", err)
	   		panic(err)
	   	}
	   	airCondition, err := strconv.ParseBool(accomodation.AirCondition)
	   	if err != nil {
	   		fmt.Println("Error parsing bool:", err)
	   		panic(err)
	   	}
	   	freeParking, err := strconv.ParseBool(accomodation.FreeParking)
	   	if err != nil {
	   		fmt.Println("Error parsing bool:", err)
	   		panic(err)
	   	}
	   	autoApproval, err := strconv.ParseBool(accomodation.AutoApproval)
	   	if err != nil {
	   		fmt.Println("Error parsing bool:", err)
	   		panic(err)
	   	} */

	//decodedImage, err := base64.StdEncoding.DecodeString(accomodation.Photos)

	return model.Accomodation{
		Name:          accomodation.Name,
		Location:      accomodation.Location,
		Wifi:          accomodation.Wifi,
		Kitchen:       accomodation.Kitchen,
		AirCondition:  accomodation.AirCondition,
		FreeParking:   accomodation.FreeParking,
		PricePerGuest: accomodation.PricePerGuest,
		AutoApproval:  accomodation.AutoApproval,
		Photos:        []byte(accomodation.Photos),
		MinGuests:     min,
		MaxGuests:     max,
		IDHost:        hostID,
	}
}

func mapSlotFromUpdateAvailability(slot *pb.UpdateAvailabilityRequest) model.Availability {

	accomodationID, err := uuid.Parse(slot.AccomodationId)
	if err != nil {
		panic(err)
	}

	price, err := strconv.ParseFloat(slot.Price, 64)
	if err != nil {
		panic(err)
	}

	return model.Availability{
		StartDate:      slot.StartDate,
		EndDate:        slot.EndDate,
		IdAccomodation: accomodationID,
		Price:          price,
	}
}

func mapAccomodation(accomodation *model.Accomodation) *pb.Accomodation {
	accomodationPb := &pb.Accomodation{
		Id:            accomodation.ID.String(),
		Name:          accomodation.Name,
		Location:      accomodation.Location,
		Wifi:          accomodation.Wifi,
		Kitchen:       accomodation.Kitchen,
		AirCondition:  accomodation.AirCondition,
		FreeParking:   accomodation.FreeParking,
		AutoApproval:  accomodation.AutoApproval,
		Photos:        accomodation.Photos,
		MinGuests:     strconv.Itoa(accomodation.MinGuests),
		MaxGuests:     strconv.Itoa(accomodation.MaxGuests),
		IdHost:        accomodation.IDHost.String(),
		PricePerGuest: accomodation.PricePerGuest,
	}
	return accomodationPb
}

func mapAvailability(availability *model.Availability) *pb.Availability {
	availabilityPb := &pb.Availability{
		Id:             availability.ID.String(),
		StartDate:      availability.StartDate,
		EndDate:        availability.EndDate,
		IdAccomodation: availability.IdAccomodation.String(),
		Price:          fmt.Sprintf("%.2f", availability.Price),
	}
	return availabilityPb
}
