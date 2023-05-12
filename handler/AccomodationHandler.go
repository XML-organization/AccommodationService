package handler

import (
	"accomodation-service/model"
	"accomodation-service/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/XML-organization/common/proto/accomodation_service"
	"github.com/google/uuid"
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
	message, err := handler.Service.AddOrUpdateAvailability(slot)
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

	accomodationID, err := uuid.Parse(in.AccommodationId)
	if err != nil {
		panic(err)
	}
	accommodation, err := handler.Service.Repo.FindByID(accomodationID)
	if err != nil {
		return &pb.AutoApprovalResponse{
			AutoApproval: accommodation.AutoApproval,
		}, err
	}

	return &pb.AutoApprovalResponse{
		AutoApproval: false,
	}, err
}

func (handler *AccomodationHandler) SearchAccomodation(w http.ResponseWriter, r *http.Request) {
	var accomodation model.AccomodationSearch
	err := json.NewDecoder(r.Body).Decode(&accomodation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//pronadji sve na toj lokaciji
	/*err = handler.Service.SearchAccByLocation(&accomodation.Location)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}*/

}
