package main

import (
	"accomodation-service/handler"
	"accomodation-service/model"
	"accomodation-service/repository"
	"accomodation-service/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	connectionStr := "postgresql://postgres:password@localhost/Airline?sslmode=disable"
	database, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.Accomodation{})
	database.Exec("SELECT 1")
	return database
}

func initRepo(database *gorm.DB) *repository.AccomodationRepository {
	return &repository.AccomodationRepository{Database: database}
}

func initServices(repo *repository.AccomodationRepository) *service.AccomodationService {
	return &service.AccomodationService{Repo: repo}
}

func initHandler(service *service.AccomodationService) *handler.AccomodationHandler {
	return &handler.AccomodationHandler{Service: service}
}
func startServer(handler *handler.AccomodationHandler) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/createAccommodation", handler.CreateAccomodation).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.0.17:5173", "http://localhost:5173", "http://192.168.137.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler1 := corsHandler.Handler(router)

	router.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent)
		})
	fmt.Println("server running ")
	log.Fatal(http.ListenAndServe(":8082", handler1))
}

func main() {
	database := initDB()
	repo := initRepo(database)
	service := initServices(repo)
	handler := initHandler(service)
	startServer(handler)
}
