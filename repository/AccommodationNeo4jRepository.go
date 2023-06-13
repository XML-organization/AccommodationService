package repository

import (
	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

type AccommodationNeo4jRepository struct {
	Session neo4j.Session
}

func NewAccommodationNeo4jRepository(driver neo4j.Driver) *AccommodationNeo4jRepository {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil
	}

	return &AccommodationNeo4jRepository{
		Session: session,
	}
}

func (repo *AccommodationNeo4jRepository) Close() {
	repo.Session.Close()
}

func (repo *AccommodationNeo4jRepository) SaveAccommodation(accommodationId string) error {
	println("Accommodation id prilikom cuvanja cvora: " + accommodationId)
	_, err := repo.Session.Run("CREATE (:Accommodation {idInPostgre: $accommodationId})",
		map[string]interface{}{
			"accommodationId": accommodationId,
		})

	if err != nil {
		println(err.Error())
		return err
	}

	return nil
}
