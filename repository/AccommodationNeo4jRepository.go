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

func (repo *AccommodationNeo4jRepository) GetAccommodationRecommentadions(userId string) ([]string, error) {
	session := repo.Session

	// Pronalaženje sličnih korisnika
	similarUsersQuery := `
		MATCH (u1:User {idInPostgre: $userId})-[:Reservation]->(:Accommodation)<-[:Reservation]-(u2:User)
		WHERE u1 <> u2
		WITH u1, u2
		MATCH (u2)-[r:Rate]->(a:Accommodation)
		WITH a, AVG(r.grade) AS averageGrade
		WHERE averageGrade >= 3
		WITH a, COUNT(*) AS ratingCount
		WHERE ratingCount >= 5
		RETURN a.idInPostgre AS accommodationId, averageGrade, ratingCount
		ORDER BY averageGrade DESC
		LIMIT 10
	`

	result, err := session.Run(similarUsersQuery, map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		return nil, err
	}

	// Sakupljanje preporučenih smeštaja
	accommodationIds := make([]string, 0)
	for result.Next() {
		record := result.Record()
		accommodationId, _ := record.Get("accommodationId")
		accommodationIds = append(accommodationIds, accommodationId.(string))
	}

	return accommodationIds, nil
}
