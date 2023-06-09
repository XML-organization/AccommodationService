package repository

import (
	"accomodation-service/model"

	neo4j "github.com/neo4j/neo4j-go-driver/neo4j"
)

type AccommodationRateNeo4jRepository struct {
	Session neo4j.Session
}

func NewAccommodationRateNeo4jRepository(driver neo4j.Driver) *AccommodationRateNeo4jRepository {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil
	}

	return &AccommodationRateNeo4jRepository{
		Session: session,
	}
}

func (repo *AccommodationRateNeo4jRepository) Close() {
	repo.Session.Close()
}

func (repo *AccommodationRateNeo4jRepository) SaveRating(rate model.HostGrade) error {
	session := repo.Session
	// Izvrši upit za kreiranje veze između korisnika i smeštaja koje je ocenio
	res, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(
			`
			MATCH (u:User {idInPostgre: $userId})
			MATCH (a:Accommodation {idInPostgre: $accommodationId})
			CREATE (u)-[:Rate {
				id: $rateId,
				grade: $grade,
				userName: $userName,
				userSurname: $userSurname,
				date: $date
			}]->(a)
			`,
			map[string]interface{}{
				"userId":          rate.UserId.String(),
				"accommodationId": rate.AccommodationId.String(),
				"rateId":          rate.ID.String(),
				"grade":           rate.Grade,
				"userName":        rate.UserName,
				"userSurname":     rate.UserSurname,
				"date":            rate.Date,
			},
		)
		if err != nil {
			return nil, err
		}
		println("USPJEŠNO SACUVANA OCJENA")
		return result, nil
	})

	if err != nil {
		println("GREŠKA PRI ČUVANJU OCJENE")
	}

	println(res)

	return nil
}

func (repo *AccommodationRateNeo4jRepository) GetAccommodationAverageRate(accommodationId string) (float64, error) {
	session := repo.Session

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		params := map[string]interface{}{
			"accommodationIdInPostgre": accommodationId,
		}

		query := `
			MATCH (a:Accommodation {idInPostgre: $accommodationIdInPostgre})<-[:Rate]-(r:User)
			RETURN avg(r.grade) AS averageRating
		`

		cursor, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}

		if cursor.Next() {
			record := cursor.Record()
			averageRating := record.GetByIndex(0).(float64)
			return averageRating, nil
		}

		return nil, nil
	})

	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, nil
	}

	return result.(float64), nil
}
