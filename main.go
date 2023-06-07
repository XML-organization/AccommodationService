package main

import (
	"accomodation-service/startup"
	cfg "accomodation-service/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
