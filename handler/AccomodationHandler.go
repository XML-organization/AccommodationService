package handler

import (
	"accomodation-service/service"
	"context"

	pb "github.com/XML-organization/common/proto/accomodation_service"
)

type AccomodationHandler struct {
	pb.UnimplementedAccomodationServiceServer
	Service *service.AccomodationService
}

func NewAccomodationHandler(service *service.AccomodationService) *AccomodationHandler {
	return &AccomodationHandler{
		Service: service,
	}
}

func (handler *AccomodationHandler) Create(ctx context.Context, request *pb.CreateRequest) (*pb.CreateResponse, error) {
	accomodation := mapAccomodationFromCreateAccomodation(request)
	message, err := handler.Service.CreateAccomodation(accomodation)
	response := pb.CreateResponse{
		Message: message.Message,
	}

	return &response, err

}
