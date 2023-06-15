package startup

import (
	"accomodation-service/handler"
	"accomodation-service/repository"
	"accomodation-service/service"
	"accomodation-service/startup/config"
	"fmt"
	"log"
	"net"

	accomodation "github.com/XML-organization/common/proto/accomodation_service"
	saga "github.com/XML-organization/common/saga/messaging"
	"github.com/XML-organization/common/saga/messaging/nats"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "accomodation_service"
)

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	neo4jDriver := server.initNeo4jDriver()
	AccomodationRepo := server.initAccomodationRepository(postgresClient)
	AccommodationNeo4jRepo := server.initAccommodationNeo4jRepository(neo4jDriver)
	AccommodationRateNeo4jRepo := server.initAccommodationRateNeo4jRepository(neo4jDriver)
	AccomodationService := server.initAccomodationService(AccomodationRepo, AccommodationNeo4jRepo, AccommodationRateNeo4jRepo)

	AcoomodationHandler := server.initAccomodationHandler(AccomodationService)

	server.startGrpcServer(AcoomodationHandler)
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetClient(
		server.config.AccomodationDBHost, server.config.AccomodationDBUser,
		server.config.AccomodationDBPass, server.config.AccomodationDBName,
		server.config.AccomodationDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initNeo4jDriver() *neo4j.Driver {
	driver, err := repository.GetNeo4jClient("bolt://accommodation_recommendation_db:7687", "neo4j", "password")

	if err != nil {
		return nil
	}

	return &driver
}

func (server *Server) initAccommodationNeo4jRepository(driver *neo4j.Driver) *repository.AccommodationNeo4jRepository {
	return repository.NewAccommodationNeo4jRepository(*driver)
}

func (server *Server) initAccommodationRateNeo4jRepository(driver *neo4j.Driver) *repository.AccommodationRateNeo4jRepository {
	return repository.NewAccommodationRateNeo4jRepository(*driver)
}

func (server *Server) initAccomodationRepository(client *gorm.DB) *repository.AccomodationRepository {
	return repository.NewAccomodationRepository(client)
}

func (server *Server) initAccomodationService(repo *repository.AccomodationRepository, neo4jRepo *repository.AccommodationNeo4jRepository, neo4jRateRepo *repository.AccommodationRateNeo4jRepository) *service.AccomodationService {
	return service.NewAccomodationService(repo, neo4jRepo, neo4jRateRepo)
}

func (server *Server) initAccomodationHandler(service *service.AccomodationService) *handler.AccomodationHandler {
	return handler.NewAccomodationHandler(service)
}

func (server *Server) startGrpcServer(accomodationHandler *handler.AccomodationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	accomodation.RegisterAccommodationServiceServer(grpcServer, accomodationHandler)
	reflection.Register(grpcServer)
	log.Println("GRPC ACCOMMODATION SERVER USPJESNO NAPRAVLJEN")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
		println("GRPC ACCOMMODATION SERVER NIJE USPJESNO NAPRAVLJEN")
	}
}
