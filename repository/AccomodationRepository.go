package repository

import (
	"accomodation-service/model"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccomodationRepository struct {
	DatabaseConnection *gorm.DB
}

func NewAccomodationRepository(db *gorm.DB) *AccomodationRepository {
	err := db.AutoMigrate(&model.Accomodation{}, &model.Availability{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return &AccomodationRepository{
		DatabaseConnection: db,
	}
}

func (repo *AccomodationRepository) CreateAccomodation(accomodation *model.Accomodation) model.RequestMessage {
	println("Id accomodationa prilikom cuvanje u postgre: " + accomodation.ID.String())
	dbResult := repo.DatabaseConnection.Save(accomodation)
	println()
	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		println(dbResult.Error.Error())
		return model.RequestMessage{
			Message: "An error occured, please try again!",
		}
	}

	return model.RequestMessage{
		Message: "Success!",
	}
}

func (repo *AccomodationRepository) UpdateAccomodation(accomodationId uuid.UUID, name string) error {
	result := repo.DatabaseConnection.Model(&model.Accomodation{}).Where("id = ?", accomodationId).Update("name", name)
	log.Println(result.RowsAffected)
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *AccomodationRepository) CreateAvailability(availability *model.Availability) model.RequestMessage {
	dbResult := repo.DatabaseConnection.Save(availability)

	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		println(dbResult.Error)
		return model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
	}

	return model.RequestMessage{
		Message: "Success!",
	}
}

func (repo *AccomodationRepository) UpdateAvailability(availability *model.Availability) error {
	result := repo.DatabaseConnection.Model(availability).Updates(availability)
	log.Println(result.RowsAffected)
	return nil
}

func (repo *AccomodationRepository) FindByID(id uuid.UUID) (model.Accomodation, error) {
	accomodation := model.Accomodation{}

	dbResult := repo.DatabaseConnection.First(&accomodation, "id = ?", id)

	if dbResult != nil {
		log.Println(dbResult.Error)
		return accomodation, dbResult.Error
	}

	return accomodation, nil
}

func (repo *AccomodationRepository) GetAllAvailabilityByIDAccomodation(availabilityID uuid.UUID) ([]model.Availability, error) {
	availabilities := []model.Availability{}
	result := repo.DatabaseConnection.Where("id_accomodation = ?", availabilityID).Find(&availabilities)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	return availabilities, nil
}

func (repo *AccomodationRepository) GetAllAccomodationByIDHost(hostID uuid.UUID) ([]model.Accomodation, error) {
	accomodations := []model.Accomodation{}
	result := repo.DatabaseConnection.Where("id_host = ?", hostID).Find(&accomodations)
	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}
	return accomodations, nil
}

func (repo *AccomodationRepository) DeleteAvailability(availabilityID uuid.UUID) error {
	result := repo.DatabaseConnection.Delete(&model.Availability{}, availabilityID)
	if result.Error != nil {
		log.Println(result.Error)
		return result.Error
	}
	return nil
}

func (repo *AccomodationRepository) FindByLocationAndNumOfGuests(location string, numOfGuests int) ([]model.Accomodation, model.RequestMessage) {
	var accommodations []model.Accomodation
	dbResult := repo.DatabaseConnection.Where("location = ? AND accomodations.min_guests <= ? AND accomodations.max_guests >= ?", location, numOfGuests, numOfGuests).Find(&accommodations)

	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		return nil, model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
	}

	return accommodations, model.RequestMessage{
		Message: "Success!",
	}
}

func (repo *AccomodationRepository) FindAllByHostId(id string) []string {
	var accommodations []model.Accomodation
	dbResult := repo.DatabaseConnection.Where("id_host = ?", id).Find(&accommodations)

	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		return []string{}
	}

	var accommodationIDs []string
	for _, accommodation := range accommodations {
		accommodationIDs = append(accommodationIDs, accommodation.ID.String())
	}

	return accommodationIDs
}
