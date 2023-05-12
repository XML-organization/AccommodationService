package repository

import (
	"accomodation-service/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccomodationRepository struct {
	DatabaseConnection *gorm.DB
}

func NewAccomodationRepository(db *gorm.DB) *AccomodationRepository {
	err := db.AutoMigrate(&model.Accomodation{}, &model.Availability{})
	if err != nil {
		return nil
	}

	return &AccomodationRepository{
		DatabaseConnection: db,
	}
}

func (repo *AccomodationRepository) CreateAccomodation(accomodation *model.Accomodation) model.RequestMessage {
	dbResult := repo.DatabaseConnection.Save(accomodation)

	if dbResult.Error != nil {
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
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *AccomodationRepository) CreateAvailability(availability model.Availability) model.RequestMessage {
	dbResult := repo.DatabaseConnection.Save(availability)

	if dbResult.Error != nil {
		println(dbResult.Error)
		return model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
	}

	return model.RequestMessage{
		Message: "Success!",
	}
}

func (repo *AccomodationRepository) UpdateAvailability(availability model.Availability) error {
	result := repo.DatabaseConnection.Model(&model.Availability{}).Where("id = ?", availability.ID).Updates(map[string]interface{}{
		"start_date": availability.StartDate,
		"end_date":   availability.EndDate,
	})
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *AccomodationRepository) FindByID(id uuid.UUID) (model.Accomodation, error) {
	accomodation := model.Accomodation{}

	dbResult := repo.DatabaseConnection.First(&accomodation, "id = ?", id)

	if dbResult != nil {
		return accomodation, dbResult.Error
	}

	return accomodation, nil
}

func (repo *AccomodationRepository) GetAllAvailabilityByIDAccomodation(availabilityID uuid.UUID) ([]model.Availability, error) {
	availabilities := []model.Availability{}
	result := repo.DatabaseConnection.Where("id_accomodation = ?", availabilityID).Find(&availabilities)
	if result.Error != nil {
		return nil, result.Error
	}
	return availabilities, nil
}

func (repo *AccomodationRepository) GetAllAccomodationByIDHost(hostID uuid.UUID) ([]model.Accomodation, error) {
	accomodations := []model.Accomodation{}
	result := repo.DatabaseConnection.Where("id_host = ?", hostID).Find(&accomodations)
	if result.Error != nil {
		return nil, result.Error
	}
	return accomodations, nil
}
