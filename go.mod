module accomodation-service

go 1.15

replace github.com/XML-organization/common => ../common

require (
	github.com/XML-organization/common v1.0.1-0.20230507161618-15f00ac411f2 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/rs/cors v1.9.0
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.24.7-0.20230306060331-85eaf9eeda11
)
