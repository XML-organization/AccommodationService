package handler

import (
	"accomodation-service/model"
	"accomodation-service/service"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	accomodation "github.com/XML-organization/common/proto/accomodation_service"
	pb "github.com/XML-organization/common/proto/accomodation_service"
	bookingServicepb "github.com/XML-organization/common/proto/booking_service"
	userServicepb "github.com/XML-organization/common/proto/user_service"
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
		log.Println(err)
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
		log.Println(err)
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
	if err != nil {
		log.Println(err)
	}
	response := pb.CreateResponse{
		Message: message.Message,
	}

	return &response, err
}

func (handler *AccomodationHandler) UpdateAvailability(ctx context.Context, request *pb.UpdateAvailabilityRequest) (*pb.UpdateAvailabilityResponse, error) {
	slot := mapSlotFromUpdateAvailability(request)
	message, err := handler.Service.AddOrUpdateAvailability(&slot)
	if err != nil {
		log.Println(err)
	}
	response := pb.UpdateAvailabilityResponse{
		Message: message.Message,
	}

	return &response, err
}

func (handler *AccomodationHandler) GetAllAccomodations(ctx context.Context, request *pb.GetAllAccomodationsRequest) (*pb.GetAllAccomodationsResponse, error) {
	println("Usao u Accomodation Service-----")
	hostID, err := uuid.Parse(request.HostId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	println("HostID u Accomodation Servicu poslije parsiranja: ", hostID.String())

	accommodations, err := handler.Service.GetAllAccomodationsByIDHost(hostID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	print("Lista koja je ucitana iz baze")
	println("Dugacka je: ", len(accommodations))
	for j, tmp := range accommodations {
		println(j, ". Accomodation ID", tmp.Name)
	}

	response := &pb.GetAllAccomodationsResponse{
		Accomodations: []*pb.Accomodation{},
	}
	for _, accomodation := range accommodations {
		current := mapAccomodation(&accomodation)
		response.Accomodations = append(response.Accomodations, current)
	}

	println("Accomodations koje vraca Accomodation Servis:")
	println("Dugacka je: ", len(response.Accomodations))

	for j, tmp := range response.Accomodations {
		println(j, ". Accomodation ID", tmp.Name)
	}

	return response, nil
}

func (handler *AccomodationHandler) GetAutoApprovalForAccommodation(ctx context.Context, in *pb.AutoApprovalRequest) (*pb.AutoApprovalResponse, error) {

	log.Println("U METODU GetAutoApprovalForAccommodation STIGAO:", in.AccommodationId)
	accomodationID, err := uuid.Parse(in.AccommodationId)
	if err != nil {
		log.Println(err)
		log.Println("ISPARSIRAO ID OVAKO: ", accomodationID.String())
		panic(err)
	}
	accommodation, err := handler.Service.Repo.FindByID(accomodationID)
	if err != nil {
		log.Println(err)
	}
	log.Println("IZ BAZE DOBAOVIO OVAJ APPROVAL: ", accommodation.AutoApproval, "I OVAJ ID", accommodation.ID.String())

	return &pb.AutoApprovalResponse{
		AutoApproval: accommodation.AutoApproval,
	}, err
}

func (handler *AccomodationHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.AccomodationSearchResponse, error) {

	searchRequest := mapAccomodationSearchFromSearchRequest(request)

	//Filtriranje prema lokaciji i broju gostiju
	accommodations, requestMessage := handler.Service.FindByLocationAndNumOfGuests(searchRequest.Location, searchRequest.NumOfGuests)
	if requestMessage.Message != "Success!" {
		log.Println("an error occurred:", requestMessage.Message)
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
			log.Println(err)
			return nil, err
		}

		for _, availability := range availabilities {

			startDate, err := time.Parse("2006-01-02", availability.StartDate)
			if err != nil {
				log.Println("Error whiile parsing date:", err)
				fmt.Println("Error whiile parsing date:", err)
				return nil, err
			}

			endDate, err := time.Parse("2006-01-02", availability.EndDate)
			if err != nil {
				log.Println("Error whiile parsing date:", err)
				fmt.Println("Error whiile parsing date:", err)
				return nil, err
			}

			if !startDate.After(start) && start.Before(endDate) {

				duration := endDate.Sub(start)
				daysDiff := int(duration.Hours() / 24)

				if daysDiff >= numOfDays {
					if accommodation.PricePerGuest {
						totalPrice = totalPrice + int(availability.Price)*numOfDays*searchRequest.NumOfGuests
						availableAccommodations = append(availableAccommodations, *mapAccomodationOnAccommodationDTO(&accommodation, totalPrice))
					} else {
						totalPrice = totalPrice + int(availability.Price)*numOfDays
						availableAccommodations = append(availableAccommodations, *mapAccomodationOnAccommodationDTO(&accommodation, totalPrice))
					}
				} else {
					if accommodation.PricePerGuest {
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
		log.Println(err)
	}
	defer conn.Close()

	bookingService := bookingServicepb.NewBookingServiceClient(conn)

	bookings, err := bookingService.GetAll(context.TODO(), &bookingServicepb.EmptyRequst{})
	if err != nil {
		log.Println(err)
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
					log.Println("Error whiile parsing date:", err)
					return nil, err
				}

				endDate, err := time.Parse("2006-01-02", booking.EndDate)
				if err != nil {
					log.Println("Error whiile parsing date:", err)
					return nil, err
				}
				log.Println("POZVAO RANGESOVERLAP ZA", startDate.String(), endDate.String(), searchRequest.StartDate.String(), searchRequest.EndDate.String())
				if rangesOverlap(startDate, endDate, searchRequest.StartDate, searchRequest.EndDate) {
					//izbaci smjestaj iz liste dostupnih
					log.Println("ovi datumi se preklapaju ", startDate.String(), endDate.String(), searchRequest.StartDate.String(), searchRequest.EndDate.String())
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
		log.Println(err)
		return nil, err
	}

	availabilities, err := handler.Service.GetAllAvailabilitiesByAccomodationID(accomodationID)
	if err != nil {
		log.Println(err)
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

func (handler *AccomodationHandler) GetAccomodations(ctx context.Context, empty *accomodation.EmptyRequst) (*pb.GetAllAccomodationsResponse, error) {

	accommodations, err := handler.Service.GetAccomodations()
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

func (handler *AccomodationHandler) GradeHost(ctx context.Context, request *pb.GradeHostRequest) (*pb.GradeHostResponse, error) {
	hostGrade := mapHostGradeFromRequest(request)
	message, err := handler.Service.GradeHost(&hostGrade)
	response := pb.GradeHostResponse{
		Message: message.Message,
	}

	accomodation, _ := handler.Service.FindByID(hostGrade.AccommodationId)
	//slanje notifikacije

	conn, err := grpc.Dial("user_service:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	userService := userServicepb.NewUserServiceClient(conn)

	saveResponse, err := userService.SaveNotification(context.TODO(), &userServicepb.SaveRequest{Id: uuid.NewString(), NotificationTime: time.Now().Format("2006-01-02 15:04:05"), Text: "Korisnik " + hostGrade.UserName + " " + hostGrade.UserSurname + " je ocijeno Vas smjestaj ocjenom: " + strconv.FormatFloat(hostGrade.Grade, 'f', -1, 64) + "  !", UserID: accomodation.IDHost.String(), Status: "0", Category: "AccommodationGraded"})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	println(saveResponse.Message)

	return &response, err
}

func (handler *AccomodationHandler) GetAccommodationRecommendations(ctx context.Context, request *pb.GetAccommodationRecommendationsRequest) (*pb.GetAccommodationRecommendationsResponse, error) {
	accommodations, err := handler.Service.GetAccommodationRecommendations(request.UserId)

	if err != nil {
		log.Println("Some error occurred when recommending an accommodation!")
		return mapAccommodationsToResponse(accommodations), err
	}

	return mapAccommodationsToResponse(accommodations), nil
}

func (handler *AccomodationHandler) GetGradesByAccomodationId(ctx context.Context, request *pb.GradesByAccomodationIdRequest) (*pb.GradesByAccomodationIdResponse, error) {
	accomodationId, err := uuid.Parse(request.AccomodationId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	grades, err := handler.Service.GetGradesByAccomodationId(accomodationId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	response := &pb.GradesByAccomodationIdResponse{
		GradesByAccomodationId: []*pb.GradeHost{},
	}
	for _, grade := range grades {
		current := mapGrades(&grade)
		response.GradesByAccomodationId = append(response.GradesByAccomodationId, current)
	}
	return response, nil
}

func (handler *AccomodationHandler) EditGrade(ctx context.Context, request *pb.EditGradeRequest) (*pb.EditGradeResponse, error) {
	gradeID, err := uuid.Parse(request.Id)
	if err != nil {
		panic(err)
	}
	gradeValue, err := strconv.ParseFloat(request.NewGrade, 64)
	if err != nil {
		log.Println("Error parsing grade:", err)
	}

	err1 := handler.Service.EditGrade(gradeID, gradeValue)
	if err1 != nil {
		log.Println(err)
		response := pb.EditGradeResponse{
			Message: "Error with edit grade!",
		}

		return &response, err1
	}
	response := pb.EditGradeResponse{
		Message: "Sucessfulu edit grade!",
	}

	return &response, nil
}

func (handler *AccomodationHandler) DeleteGrade(ctx context.Context, request *pb.DeleteGradeRequest) (*pb.DeleteGradeResponse, error) {
	gradeID, err := uuid.Parse(request.Id)
	if err != nil {
		panic(err)
	}

	err1 := handler.Service.DeleteGrade(gradeID)
	if err1 != nil {
		log.Println(err)
		response := pb.DeleteGradeResponse{
			Message: "Error with delete grade!",
		}

		return &response, err1
	}
	response := pb.DeleteGradeResponse{
		Message: "Sucessfulu delete grade!",
	}

	return &response, nil
}
