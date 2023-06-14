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

func (repo *AccommodationNeo4jRepository) FindSimilarUsers(userID string) ([]string, error) {
	session := repo.Session

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		params := map[string]interface{}{
			"userID": userID,
		}

		query := `
			MATCH (u:User {idInPostgre: $userID})-[:Reservation]->(a:Accommodation)
			WITH a, u
			MATCH (u2:User)-[:Reservation]->(a)
			WHERE u2 <> u
			RETURN DISTINCT u2.idInPostgre AS similarUserID
		`

		cursor, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		var similarUserIDs []string
		for cursor.Next() {
			record := cursor.Record()
			similarUserID := record.GetByIndex(0).(string)
			similarUserIDs = append(similarUserIDs, similarUserID)
		}

		return similarUserIDs, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

func (repo *AccommodationNeo4jRepository) FindRecommendedAccommodations(similarUserIDs []string) ([]string, error) {
	session := repo.Session

	result, err := session.Run(`
		MATCH (u:User)-[:Reservation]->(a:Accommodation)<-[r:Rate]-(rUser:User)
		WHERE u.idInPostgre IN $similarUserIDs and r.grade > 3
		WITH a, COLLECT(DISTINCT u.idInPostgre) AS userIDs
		WHERE SIZE(userIDs) < SIZE($similarUserIDs)
		RETURN DISTINCT a.idInPostgre AS recommendedAccommodationID

	`, map[string]interface{}{
		"similarUserIDs": similarUserIDs,
	})

	if err != nil {
		return nil, err
	}

	var recommendedAccommodationIDs []string
	for result.Next() {
		record := result.Record()
		recommendedAccommodationIDValue, found := record.Get("recommendedAccommodationID")
		if found {
			recommendedAccommodationID, ok := recommendedAccommodationIDValue.(string)
			if ok {
				recommendedAccommodationIDs = append(recommendedAccommodationIDs, recommendedAccommodationID)
			}
		}
	}

	return recommendedAccommodationIDs, nil
}

func (repo *AccommodationNeo4jRepository) FilterAccommodations(recommendedAccommodationIDs []string) ([]string, error) {
	session := repo.Session

	result, err := session.Run(`
		MATCH (a:Accommodation)<-[r:Rate]-(:User)
		WHERE a.idInPostgre IN $recommendedAccommodationIDs
		WITH a, SUM(CASE WHEN r.grade < 3 AND r.date >= datetime() - duration({months: 3}) THEN 1 ELSE 0 END) AS lowRatingsCount
		WHERE lowRatingsCount <= 5
		RETURN a.idInPostgre AS filteredAccommodationID
	`, map[string]interface{}{
		"recommendedAccommodationIDs": recommendedAccommodationIDs,
	})

	if err != nil {
		return nil, err
	}

	var filteredAccommodationIDs []string
	for result.Next() {
		record := result.Record()
		filteredAccommodationIDValue, found := record.Get("filteredAccommodationID")
		if found {
			filteredAccommodationID, ok := filteredAccommodationIDValue.(string)
			if ok {
				filteredAccommodationIDs = append(filteredAccommodationIDs, filteredAccommodationID)
			}
		}
	}

	return filteredAccommodationIDs, nil
}

func (repo *AccommodationNeo4jRepository) RankAccommodations(filteredAccommodationIDs []string) ([]string, error) {
	session := repo.Session

	result, err := session.Run(`
		MATCH (a:Accommodation)<-[:Rate]-(r:User)
		WHERE a.idInPostgre IN $filteredAccommodationIDs
		WITH a, avg(r.grade) AS overallRating
		RETURN a.idInPostgre AS rankedAccommodationID
		ORDER BY overallRating DESC
	`, map[string]interface{}{
		"filteredAccommodationIDs": filteredAccommodationIDs,
	})

	if err != nil {
		return nil, err
	}

	var rankedAccommodationIDs []string
	for result.Next() {
		record := result.Record()
		rankedAccommodationIDValue, found := record.Get("rankedAccommodationID")
		if found {
			rankedAccommodationID, ok := rankedAccommodationIDValue.(string)
			if ok {
				rankedAccommodationIDs = append(rankedAccommodationIDs, rankedAccommodationID)
			}
		}
	}

	return rankedAccommodationIDs, nil
}
