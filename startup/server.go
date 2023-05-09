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
	"google.golang.org/grpc"
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
	AccomodationRepo := server.initAccomodationRepository(postgresClient)
	AccomodationService := server.initAccomodationService(AccomodationRepo)

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

func (server *Server) initAccomodationRepository(client *gorm.DB) *repository.AccomodationRepository {
	return repository.NewAccomodationRepository(client)
}

func (server *Server) initAccomodationService(repo *repository.AccomodationRepository) *service.AccomodationService {
	return service.NewAccomodationService(repo)
}

func (server *Server) initAccomodationHandler(service *service.AccomodationService) *handler.AccomodationHandler {
	return handler.NewAccomodationHandler(service)
}

func (server *Server) startGrpcServer(accomodationHandler *handler.AccomodationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	accomodation.RegisterAccomodationServiceServer(grpcServer, accomodationHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
