package service

import (
	"accomodation-service/model"
	"accomodation-service/repository"
)

type AccomodationService struct {
	Repo *repository.AccomodationRepository
}

func (service *AccomodationService) CreateAccomodation(accomodation *model.Accomodation) error {

	err := service.Repo.CreateAccomodation(accomodation)
	if err != nil {
		return err
	}
	return nil
}
