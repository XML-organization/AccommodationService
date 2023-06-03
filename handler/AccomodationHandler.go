package handler

import (
	"accomodation-service/model"
	"accomodation-service/service"
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/XML-organization/common/proto/accomodation_service"
	bookingServicepb "github.com/XML-organization/common/proto/booking_service"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccomodationHandler struct {
	pb.UnimplementedAccommodationServiceServer
	Service *service.AccomodationService
}

func NewAccomodationHandler(service *service.AccomodationService) *AccomodationHandler {
	return &AccomodationHandler{
		Service: service,
	}
}

func (handler *AccomodationHandler) CheckIfGuestHasReservationInPast(ctx context.Context, request *pb.CheckIfGuestHasReservationInPastRequest) (*pb.CheckIfGuestHasReservationInPastResponse, error) {

	ids := handler.Service.FindAllAccomodationIDsByHostId(request.HostId)

	conn, err := grpc.Dial("booking-service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	bookingService := bookingServicepb.NewBookingServiceClient(conn)

	println(ids[0])

	hasResevation, err := bookingService.GuestHasReservationInPast(context.TODO(), &bookingServicepb.GuestHasReservationInPastRequest{AccomodationsId: ids, GuestId: request.GuestId})
	if err != nil {
		println(err.Error())
		return nil, err
	}

	retValue := false
	if hasResevation.Message == "Have" {
		retValue = true
	}

	return &pb.CheckIfGuestHasReservationInPastResponse{HasReservation: retValue}, nil
}

func (handler *AccomodationHandler) GetOneAccomodation(ctx context.Context, request *pb.GetOneAccomodationRequest) (*pb.GetOneAccomodationResponse, error) {
	accomodationID, err := uuid.Parse(request.AccomodationId)
	if err != nil {
		return &pb.GetOneAccomodationResponse{}, err
	}

	accomodation, _ := handler.Service.FindByID(accomodationID)

	return &pb.GetOneAccomodationResponse{
		Accomodation: mapAccomodation(&accomodation),
	}, nil
}

func (handler *AccomodationHandler) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	accomodation := mapAccomodationFromCreateAccomodation(request)
	message, err := handler.Service.CreateAccomodation(&accomodation)
	response := pb.CreateResponse{
		Message: message.Message,
	}

	return &response, err
}

func (handler *AccomodationHandler) UpdateAvailability(ctx context.Context, request *pb.UpdateAvailabilityRequest) (*pb.UpdateAvailabilityResponse, error) {
	slot := mapSlotFromUpdateAvailability(request)
	message, err := handler.Service.AddOrUpdateAvailability(&slot)
	response := pb.UpdateAvailabilityResponse{
		Message: message.Message,
	}

	return &response, err
}

func (handler *AccomodationHandler) GetAllAccomodations(ctx context.Context, request *pb.GetAllAccomodationsRequest) (*pb.GetAllAccomodationsResponse, error) {
	hostID, err := uuid.Parse(request.HostId)
	if err != nil {
		return nil, err
	}

	accommodations, err := handler.Service.GetAllAccomodationsByIDHost(hostID)
	if err != nil {
		return nil, err
	}

	response := &pb.GetAllAccomodationsResponse{
		Accomodations: []*pb.Accomodation{},
	}
	for _, accomodation := range accommodations {
		current := mapAccomodation(&accomodation)
		response.Accomodations = append(response.Accomodations, current)
	}
	return response, nil
}

func (handler *AccomodationHandler) GetAutoApprovalForAccommodation(ctx context.Context, in *pb.AutoApprovalRequest) (*pb.AutoApprovalResponse, error) {

	println("U METODU GetAutoApprovalForAccommodation STIGAO:", in.AccommodationId)
	accomodationID, err := uuid.Parse(in.AccommodationId)
	if err != nil {
		println("ISPARSIRAO ID OVAKO: ", accomodationID.String())
		panic(err)
	}
	accommodation, err := handler.Service.Repo.FindByID(accomodationID)
	println("IZ BAZE DOBAOVIO OVAJ APPROVAL: ", accommodation.AutoApproval, "I OVAJ ID", accommodation.ID.String())

	return &pb.AutoApprovalResponse{
		AutoApproval: accommodation.AutoApproval,
	}, err
}

func (handler *AccomodationHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.AccomodationSearchResponse, error) {

	searchRequest := mapAccomodationSearchFromSearchRequest(request)

	//Filtriranje prema lokaciji i broju gostiju
	accommodations, requestMessage := handler.Service.FindByLocationAndNumOfGuests(searchRequest.Location, searchRequest.NumOfGuests)
	if requestMessage.Message != "Success!" {
		return nil, fmt.Errorf("an error occurred: %s", requestMessage.Message)
	}

	//Provjera dostupnosti objekta i cijene u vremenskom intervalu

	availableAccommodations := []model.AccomodationDTO{}

	for _, accommodation := range accommodations {

		start := searchRequest.StartDate
		end := searchRequest.EndDate
		numOfDays := int(end.Sub(start).Hours() / 24)
		totalPrice := 0

		availabilities, err := handler.Service.GetAllAvailabilitiesByAccomodationID(accommodation.ID)
		if err != nil {
			return nil, err
		}

		for _, availability := range availabilities {

			startDate, err := time.Parse("2006-01-02", availability.StartDate)
			if err != nil {
				fmt.Println("Error whiile parsing date:", err)
				return nil, err
			}

			endDate, err := time.Parse("2006-01-02", availability.EndDate)
			if err != nil {
				fmt.Println("Error whiile parsing date:", err)
				return nil, err
			}

			if !startDate.After(start) && start.Before(endDate) {

				duration := endDate.Sub(start)
				daysDiff := int(duration.Hours() / 24)

				if daysDiff >= numOfDays {
					if accommodation.PricePerGuest == true {
						totalPrice = totalPrice + int(availability.Price)*numOfDays*searchRequest.NumOfGuests
						availableAccommodations = append(availableAccommodations, *mapAccomodationOnAccommodationDTO(&accommodation, totalPrice))
					} else {
						totalPrice = totalPrice + int(availability.Price)*numOfDays
						availableAccommodations = append(availableAccommodations, *mapAccomodationOnAccommodationDTO(&accommodation, totalPrice))
					}
				} else {
					if accommodation.PricePerGuest == true {
						totalPrice = totalPrice + int(availability.Price)*daysDiff*searchRequest.NumOfGuests
						numOfDays = numOfDays - daysDiff
						start = endDate
					} else {
						totalPrice = totalPrice + int(availability.Price)*daysDiff
						numOfDays = numOfDays - daysDiff
						start = endDate
					}
				}

			} else {
				continue
			}
		}

	}

	//rpc GetAllBookings
	conn, err := grpc.Dial("booking-service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	bookingService := bookingServicepb.NewBookingServiceClient(conn)

	bookings, err := bookingService.GetAll(context.TODO(), &bookingServicepb.EmptyRequst{})
	if err != nil {
		println(err.Error())
		return nil, err
	}

	//Provjera da li je smjestaj vec rezervisan u navedenom periodu

	for _, accomodation := range availableAccommodations {
		println("------------BOOKINGS----------")
		for _, booking := range bookings.Bookings {
			println(booking.Status)
			println(booking.Id)

			if booking.Status == "CONFIRMED" && booking.AccomodationID == accomodation.ID.String() {

				startDate, err := time.Parse("2006-01-02", booking.StartDate)
				if err != nil {
					fmt.Println("Error whiile parsing date:", err)
					return nil, err
				}

				endDate, err := time.Parse("2006-01-02", booking.EndDate)
				if err != nil {
					fmt.Println("Error whiile parsing date:", err)
					return nil, err
				}
				println("POZVAO RANGESOVERLAP ZA", startDate.String(), endDate.String(), searchRequest.StartDate.String(), searchRequest.EndDate.String())
				if rangesOverlap(startDate, endDate, searchRequest.StartDate, searchRequest.EndDate) {
					//izbaci smjestaj iz liste dostupnih
					println("ovi datumi se preklapaju ", startDate.String(), endDate.String(), searchRequest.StartDate.String(), searchRequest.EndDate.String())
					removeAccommodationFromList(&availableAccommodations, &accomodation)
				}
			}
		}
	}

	response := pb.AccomodationSearchResponse{
		AccommodationsDTO: []*pb.AccomodationDTO{},
	}

	for _, accommodation := range availableAccommodations {
		proto := mapAccomodationDTOToAccommodationSearchResponse(&accommodation)
		response.AccommodationsDTO = append(response.AccommodationsDTO, proto)
	}

	return &response, nil

}

func rangesOverlap(start1, end1, start2, end2 time.Time) bool {
	return !(end1.Before(start2) || end2.Before(start1) || start1.Equal(end2) || start2.Equal(end1))
}

func removeAccommodationFromList(accommodations *[]model.AccomodationDTO, accommodation *model.AccomodationDTO) {
	println("pozvao sam metodu ukloni smjestaj iz liste dostupnih")
	for i, acc := range *accommodations {
		if acc.ID == accommodation.ID {
			// Pronađen objekat, uklanjanje iz liste
			println("OBRISAO")
			*accommodations = append((*accommodations)[:i], (*accommodations)[i+1:]...)
			break
		}
	}
}

func (handler *AccomodationHandler) GetAllAvailability(ctx context.Context, request *pb.GetAllAvailabilityRequest) (*pb.GetAllAvailabilityResponse, error) {
	accomodationID, err := uuid.Parse(request.AccomodationId)
	if err != nil {
		return nil, err
	}

	availabilities, err := handler.Service.GetAllAvailabilitiesByAccomodationID(accomodationID)
	if err != nil {
		return nil, err
	}

	response := &pb.GetAllAvailabilityResponse{
		Availabilities: []*pb.Availability{},
	}
	for _, availability := range availabilities {
		current := mapAvailability(&availability)
		response.Availabilities = append(response.Availabilities, current)
	}
	return response, nil
}
