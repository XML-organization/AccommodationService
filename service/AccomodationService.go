package service

import (
	"accomodation-service/model"
	"accomodation-service/repository"
)

type AccomodationService struct {
	Repo *repository.AccomodationRepository
}

func NewAccomodationService(repo *repository.AccomodationRepository) *AccomodationService {
	return &AccomodationService{
		Repo: repo,
	}
}

func (service *AccomodationService) CreateAccomodation(accomodation model.Accomodation) (model.RequestMessage, error) {

	response := model.RequestMessage{
		Message: service.Repo.CreateAccomodation(accomodation).Message,
	}

	return response, nil
}
