package service

import (
	"accomodation-service/model"
	"accomodation-service/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type AccomodationService struct {
	Repo *repository.AccomodationRepository
}

func NewAccomodationService(repo *repository.AccomodationRepository) *AccomodationService {
	return &AccomodationService{
		Repo: repo,
	}
}

func (service *AccomodationService) CreateAccomodation(accomodation *model.Accomodation) (model.RequestMessage, error) {
	println("acc servis")
	println(strconv.FormatBool(accomodation.AutoApproval))
	println(accomodation.Name)
	println(accomodation.Photos)
	println(accomodation.IDHost.String())
	response := model.RequestMessage{
		Message: service.Repo.CreateAccomodation(accomodation).Message,
	}

	return response, nil
}

func (service *AccomodationService) AddOrUpdateAvailability(availability model.Availability) (model.RequestMessage, error) {
	existingAvailabilities, err := service.Repo.GetAllAvailabilityByIDAccomodation(availability.IdAccomodation)
	if err != nil {
		return model.RequestMessage{}, err
	}
	startDate, s1 := time.Parse("2006-01-02", availability.StartDate)

	endDate, e1 := time.Parse("2006-01-02", availability.EndDate)

	for _, existingAvailability := range existingAvailabilities {
		EAstartDate, s2 := time.Parse("2006-01-02", existingAvailability.StartDate)

		EAendDate, e2 := time.Parse("2006-01-02", existingAvailability.EndDate)

		// Provera preklapanja sa postojećim dostupnim terminom
		if isTimeOverlap(existingAvailability, availability) {
			// Postoji preklapanje vremena
			if s1 != nil && s2 != nil && e1 != nil && e2 != nil {
				if EAstartDate.Before(startDate) {
					// Preklapanje se događa na početku postojećeg termina
					EAendDate = startDate
				} else if EAendDate.After(endDate) {
					// Preklapanje se događa na kraju postojećeg termina
					EAstartDate = endDate
				}

				// Sačuvaj promene u repozitorijumu
				err = service.Repo.UpdateAvailability(existingAvailability)
				if err != nil {
					// Greška pri ažuriranju postojećeg dostupnog termina
					return model.RequestMessage{}, err
				}
			}
		}
	}

	response := model.RequestMessage{
		Message: service.Repo.CreateAvailability(availability).Message,
	}

	return response, nil
}

// Pomoćna funkcija za proveru preklapanja vremena
func isTimeOverlap(availability1, availability2 model.Availability) bool {
	startDate1, s1 := time.Parse("2006-01-02", availability1.StartDate)
	endDate1, e1 := time.Parse("2006-01-02", availability1.EndDate)
	startDate2, s2 := time.Parse("2006-01-02", availability2.StartDate)
	endDate2, e2 := time.Parse("2006-01-02", availability2.EndDate)
	if s1 != nil && s2 != nil && e1 != nil && e2 != nil {
		return startDate1.Before(endDate2) && startDate2.Before(endDate1)
	} else {
		return false
	}
}

func (service *AccomodationService) GetAllAccomodationsByIDHost(hostID uuid.UUID) ([]model.Accomodation, error) {
	accomodations, err := service.Repo.GetAllAccomodationByIDHost(hostID)
	if err != nil {
		return nil, err
	}
	return accomodations, nil
}
