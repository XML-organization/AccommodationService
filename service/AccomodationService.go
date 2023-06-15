package service

import (
	"accomodation-service/model"
	"accomodation-service/repository"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type AccomodationService struct {
	Repo          *repository.AccomodationRepository
	Neo4jRepo     *repository.AccommodationNeo4jRepository
	Neo4jRateRepo *repository.AccommodationRateNeo4jRepository
}

func NewAccomodationService(repo *repository.AccomodationRepository, neo4jRepo *repository.AccommodationNeo4jRepository, rateRepo *repository.AccommodationRateNeo4jRepository) *AccomodationService {
	return &AccomodationService{
		Repo:          repo,
		Neo4jRepo:     neo4jRepo,
		Neo4jRateRepo: rateRepo,
	}
}

func (service *AccomodationService) FindAllAccomodationIDsByHostId(id string) []string {

	return service.Repo.FindAllByHostId(id)
}

func (service *AccomodationService) CreateAccomodation(accomodation *model.Accomodation) (model.RequestMessage, error) {

	accomodation.ID = uuid.New()
	service.Neo4jRepo.SaveAccommodation(accomodation.ID.String())
	response := model.RequestMessage{
		Message: service.Repo.CreateAccomodation(accomodation).Message,
	}

	return response, nil
}

func (service *AccomodationService) AddOrUpdateAvailability(availability *model.Availability) (model.RequestMessage, error) {
	existingAvailabilities, err := service.Repo.GetAllAvailabilityByIDAccomodation(availability.IdAccomodation)
	if err != nil {
		log.Println(err)
		return model.RequestMessage{}, err
	}

	startDate, err := time.Parse("2006-01-02", availability.StartDate)
	if err != nil {
		log.Println(err)
		return model.RequestMessage{}, err
	}

	endDate, err := time.Parse("2006-01-02", availability.EndDate)
	if err != nil {
		log.Println(err)
		return model.RequestMessage{}, err
	}

	fmt.Println("novo vreme pocetak: " + availability.StartDate)
	fmt.Println("novo vreme kraj: " + availability.EndDate)

	for _, existingAvailability := range existingAvailabilities {
		EAstartDate, err := time.Parse("2006-01-02", existingAvailability.StartDate)
		if err != nil {
			log.Println(err)
			return model.RequestMessage{}, err
		}

		EAendDate, err := time.Parse("2006-01-02", existingAvailability.EndDate)
		if err != nil {
			log.Println(err)
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
					log.Println(err)
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
					log.Println(err)
					return model.RequestMessage{}, err
				}
			}

			// Sačuvaj promene u repozitorijumu
			err = service.Repo.UpdateAvailability(&existingAvailability)
			if err != nil {
				log.Println(err)
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
		log.Println(err)
		return nil, err
	}
	return accomodations, nil
}

func (service *AccomodationService) GetAllAvailabilitiesByAccomodationID(accomodationID uuid.UUID) ([]model.Availability, error) {
	availabilities, err := service.Repo.GetAllAvailabilityByIDAccomodation(accomodationID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return availabilities, nil
}

func (service *AccomodationService) FindByLocationAndNumOfGuests(location string, numOfGuests int) ([]model.Accomodation, model.RequestMessage) {
	accommodations, err := service.Repo.FindByLocationAndNumOfGuests(location, numOfGuests)
	if err.Message != "Success!" {
		log.Println(err)
		return nil, model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
	}
	return accommodations, err
}

func (service *AccomodationService) FindByID(id uuid.UUID) (model.Accomodation, model.RequestMessage) {
	accommodations, err := service.Repo.FindByID(id)
	if err != nil {
		log.Println(err)
		return model.Accomodation{}, model.RequestMessage{
			Message: "Accomodation not found!",
		}
	}
	return accommodations, model.RequestMessage{
		Message: "Successfully!",
	}
}

func (service *AccomodationService) GetAccomodations() ([]model.Accomodation, error) {
	accomodations, err := service.Repo.GetAccomodations()
	if err != nil {
		return nil, err
	}
	return accomodations, nil
}

func (service *AccomodationService) GradeHost(hostGrade *model.HostGrade) (model.RequestMessage, error) {
	log.Println("Call function GradeHost")

	service.Neo4jRateRepo.SaveRating(*hostGrade)

	response := model.RequestMessage{
		Message: service.Repo.GradeHost(hostGrade).Message,
	}

	return response, nil
}

func (service *AccomodationService) GetAccommodationRecommendations(userId string) ([]model.Accomodation, error) {
	log.Println("Call function GetAccommodationRecommendations")

	accommodationIds, err := service.Neo4jRepo.GetAccommodationRecommentadions(userId)

	if err != nil {
		return []model.Accomodation{}, err
	}

	log.Printf(accommodationIds[0])

	return service.Repo.FindAccommodationsByIds(accommodationIds)
}

func (service *AccomodationService) GetGradesByAccomodationId(accomodationId uuid.UUID) ([]model.HostGrade, error) {
	grades, err := service.Repo.GetGradesByAccomodationId(accomodationId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return grades, err

}

func (service *AccomodationService) EditGrade(gradeId uuid.UUID, newGrade float64) error {
	err := service.Repo.EditGrade(gradeId, newGrade)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccomodationService) DeleteGrade(gradeId uuid.UUID) error {
	err := service.Repo.DeleteGrade(gradeId)
	if err != nil {
		return err
	}
	return nil
}
