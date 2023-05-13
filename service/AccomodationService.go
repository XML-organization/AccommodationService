package service

import (
	"accomodation-service/model"
	"accomodation-service/repository"
	"fmt"
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

func (service *AccomodationService) AddOrUpdateAvailability(availability *model.Availability) (model.RequestMessage, error) {
	existingAvailabilities, err := service.Repo.GetAllAvailabilityByIDAccomodation(availability.IdAccomodation)
	if err != nil {
		return model.RequestMessage{}, err
	}

	startDate, err := time.Parse("2006-01-02", availability.StartDate)
	if err != nil {
		return model.RequestMessage{}, err
	}

	endDate, err := time.Parse("2006-01-02", availability.EndDate)
	if err != nil {
		return model.RequestMessage{}, err
	}

	fmt.Println("novo vreme pocetak: " + availability.StartDate)
	fmt.Println("novo vreme kraj: " + availability.EndDate)

	for _, existingAvailability := range existingAvailabilities {
		EAstartDate, err := time.Parse("2006-01-02", existingAvailability.StartDate)
		if err != nil {
			return model.RequestMessage{}, err
		}

		EAendDate, err := time.Parse("2006-01-02", existingAvailability.EndDate)
		if err != nil {
			return model.RequestMessage{}, err
		}

		// Provera preklapanja sa postojećim dostupnim terminom
		if startDate.Before(EAendDate) && EAstartDate.Before(endDate) {
			// Postoji preklapanje vremena
			fmt.Println("preklapa se")
			if EAstartDate.Before(startDate) && EAendDate.Before(endDate) {
				// Preklapanje: početak postojećeg termina je pre početka novog termina, a kraj postojećeg termina je pre kraja novog termina
				EAendDate = startDate
				existingAvailability.EndDate = EAendDate.Format("2006-01-02")
			} else if EAstartDate.After(startDate) && EAendDate.After(endDate) {
				// Preklapanje: početak postojećeg termina je posle početka novog termina, a kraj postojećeg termina je posle kraja novog termina
				EAstartDate = endDate
				existingAvailability.StartDate = EAstartDate.Format("2006-01-02")
			} else if EAstartDate.After(startDate) && EAendDate.Before(endDate) {
				// Preklapanje: početak postojećeg termina je posle početka novog termina, a kraj postojećeg termina je pre kraja novog termina
				// Obrisi postojeci termin
				err = service.Repo.DeleteAvailability(existingAvailability.ID)
				if err != nil {
					return model.RequestMessage{}, err
				}
			} else if EAstartDate.Before(startDate) && EAendDate.After(endDate) {
				// Preklapanje: početak postojećeg termina je pre početka novog termina, a kraj postojećeg termina je posle kraja novog termina
				// Podeli postojeci termin na dva dela
				newAvailability1 := model.Availability{
					StartDate: existingAvailability.StartDate,
					EndDate:   availability.StartDate,
					// Ostali atributi preuzeti iz existingAvailability
				}
				newAvailability2 := model.Availability{
					StartDate: availability.EndDate,
					EndDate:   existingAvailability.EndDate,
					// Ostali atributi preuzeti iz existingAvailability
				}

				// Sacuvaj nove termine u repozitorijumu
				service.Repo.CreateAvailability(&newAvailability1)

				service.Repo.CreateAvailability(&newAvailability2)

				// Obrisi postojeci termin
				err = service.Repo.DeleteAvailability(existingAvailability.ID)
				if err != nil {
					return model.RequestMessage{}, err
				}
			}

			// Sačuvaj promene u repozitorijumu
			err = service.Repo.UpdateAvailability(&existingAvailability)
			if err != nil {
				// Greška pri ažuriranju postojećeg dostupnog termina
				return model.RequestMessage{}, err
			}
		}

	}

	response := model.RequestMessage{
		Message: service.Repo.CreateAvailability(availability).Message,
	}

	return response, nil
}

func (service *AccomodationService) GetAllAccomodationsByIDHost(hostID uuid.UUID) ([]model.Accomodation, error) {
	accomodations, err := service.Repo.GetAllAccomodationByIDHost(hostID)
	if err != nil {
		return nil, err
	}
	return accomodations, nil
}

func (service *AccomodationService) GetAllAvailabilitiesByAccomodationID(accomodationID uuid.UUID) ([]model.Availability, error) {
	availabilities, err := service.Repo.GetAllAvailabilityByIDAccomodation(accomodationID)
	if err != nil {
		return nil, err
	}
	return availabilities, nil
}
