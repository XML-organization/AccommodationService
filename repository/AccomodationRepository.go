package repository

import (
	"accomodation-service/model"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccomodationRepository struct {
	Database *gorm.DB
}

func (repo *AccomodationRepository) CreateAccomodation(accomodation *model.Accomodation) error {
	result := repo.Database.Create(accomodation)
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *AccomodationRepository) UpdateAccomodation(accomodationId uuid.UUID, name string) error {
	result := repo.Database.Model(&model.Accomodation{}).Where("id = ?", accomodationId).Update("name", name)
	fmt.Println(result.RowsAffected)
	return nil
}
