package handler

import (
	"accomodation-service/model"
	"fmt"
	"log"
	"strconv"
	"time"

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
		log.Println(err)
		panic(err)
	}
	min, err := strconv.Atoi(accomodation.MinGuests)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	max, err := strconv.Atoi(accomodation.MaxGuests)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		panic(err)
	}

	price, err := strconv.ParseFloat(slot.Price, 64)
	if err != nil {
		log.Println(err)
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

func mapAccomodationOnAccommodationDTO(accomodation *model.Accomodation, totalPrice int) *model.AccomodationDTO {
	println("Id smjestaja prilikom search: " + accomodation.ID.String())
	accomodationDTO := &model.AccomodationDTO{
		ID:            accomodation.ID,
		Name:          accomodation.Name,
		Location:      accomodation.Location,
		Wifi:          accomodation.Wifi,
		Kitchen:       accomodation.Kitchen,
		AirCondition:  accomodation.AirCondition,
		FreeParking:   accomodation.FreeParking,
		AutoApproval:  accomodation.AutoApproval,
		Photos:        accomodation.Photos,
		MinGuests:     accomodation.MinGuests,
		MaxGuests:     accomodation.MaxGuests,
		IDHost:        accomodation.IDHost,
		PricePerGuest: accomodation.PricePerGuest,
		TotalPrice:    totalPrice,
	}
	return accomodationDTO
}

func mapAccomodationDTOToAccommodationSearchResponse(accomodation *model.AccomodationDTO) *pb.AccomodationDTO {
	println("Accommodation id u search responsu: " + accomodation.ID.String())
	accomodationPb := &pb.AccomodationDTO{
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
		TotalPrice:    strconv.Itoa(accomodation.TotalPrice),
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

func mapAccomodationSearchFromSearchRequest(search *pb.SearchRequest) model.AccomodationSearch {

	num, err := strconv.Atoi(search.NumOfGuests)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, search.StartDate)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	endDate, err := time.Parse(layout, search.EndDate)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	return model.AccomodationSearch{
		Location:    search.Location,
		StartDate:   startDate,
		EndDate:     endDate,
		NumOfGuests: num,
	}
}

func mapAccomodationSearchToSearchRequest(accomodationSearch *model.AccomodationSearch) *pb.SearchRequest {
	guestNumber := strconv.Itoa(accomodationSearch.NumOfGuests)
	startDate := accomodationSearch.StartDate.Format("2006-01-02")
	endDate := accomodationSearch.EndDate.Format("2006-01-02")

	return &pb.SearchRequest{
		Location:    accomodationSearch.Location,
		StartDate:   startDate,
		EndDate:     endDate,
		NumOfGuests: guestNumber,
	}
}

func mapHostGradeFromRequest(grade *pb.GradeHostRequest) model.HostGrade {

	accommodationID, err := uuid.Parse(grade.AccomodationId)
	if err != nil {
		panic(err)
	}
	userId, err := uuid.Parse(grade.UserId)
	if err != nil {
		panic(err)
	}
	gradeValue, err := strconv.ParseFloat(grade.Grade, 64)
	if err != nil {
		log.Println("Error parsing grade:", err)
	}

	layout := "2006-01-02"
	t := time.Now()
	formattedDate, _ := time.Parse(layout, t.Format(layout))

	id := uuid.New()

	return model.HostGrade{
		ID:              id,
		AccommodationId: accommodationID,
		UserId:          userId,
		UserName:        grade.UserName,
		UserSurname:     grade.UserSurname,
		Grade:           gradeValue,
		Date:            formattedDate,
	}
}

func mapAccommodationsToResponse(accommodations []model.Accomodation) *pb.GetAccommodationRecommendationsResponse {
	response := &pb.GetAccommodationRecommendationsResponse{
		Accomodations: make([]*pb.Accomodation, len(accommodations)),
	}

	for i, a := range accommodations {
		response.Accomodations[i] = &pb.Accomodation{
			Id:            a.ID.String(),
			Name:          a.Name,
			Location:      a.Location,
			Wifi:          a.Wifi,
			Kitchen:       a.Kitchen,
			AirCondition:  a.AirCondition,
			FreeParking:   a.FreeParking,
			AutoApproval:  a.AutoApproval,
			Photos:        a.Photos,
			MinGuests:     strconv.Itoa(a.MinGuests),
			MaxGuests:     strconv.Itoa(a.MaxGuests),
			IdHost:        a.IDHost.String(),
			PricePerGuest: a.PricePerGuest,
		}
	}

	return response
}
