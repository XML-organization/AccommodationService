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
	err := db.AutoMigrate(&model.Accomodation{})
	if err != nil {
		return nil
	}

	return &AccomodationRepository{
		DatabaseConnection: db,
	}
}

func (repo *AccomodationRepository) CreateAccomodation(accomodation model.Accomodation) model.RequestMessage {
	dbResult := repo.DatabaseConnection.Save(accomodation)

	if dbResult.Error != nil {
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
