package main

import (
	"accomodation-service/startup"
	cfg "accomodation-service/startup/config"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	config := cfg.NewConfig()
	log.Println("Starting server Accomodation Service...")
	server := startup.NewServer(config)
	server.Start()
}
