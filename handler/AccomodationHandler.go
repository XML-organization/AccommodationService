package handler

import (
	"accomodation-service/service"
	"context"

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
