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

func (handler *AccomodationHandler) CreateAccomodation(w http.ResponseWriter, r *http.Request) {
	var accomodation model.Accomodation
	err := json.NewDecoder(r.Body).Decode(&accomodation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.Service.CreateAccomodation(&accomodation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
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
