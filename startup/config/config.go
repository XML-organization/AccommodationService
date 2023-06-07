package config

import "os"

type Config struct {
	Port                     string
	AccomodationDBHost       string
	AccomodationDBPort       string
	AccomodationDBName       string
	AccomodationDBUser       string
	AccomodationDBPass       string
	NatsHost                 string
	NatsPort                 string
	NatsUser                 string
	NatsPass                 string
	CreateUserCommandSubject string
	CreateUserReplySubject   string
}

func NewConfig() *Config {
	return &Config{
		Port:                     os.Getenv("ACCOMODATION_SERVICE_PORT"),
		AccomodationDBHost:       os.Getenv("ACCOMODATION_DB_HOST"),
		AccomodationDBPort:       os.Getenv("ACCOMODATION_DB_PORT"),
		AccomodationDBName:       os.Getenv("ACCOMODATION_DB_NAME"),
		AccomodationDBUser:       os.Getenv("ACCOMODATION_DB_USER"),
		AccomodationDBPass:       os.Getenv("ACCOMODATION_DB_PASS"),
		NatsHost:                 os.Getenv("NATS_HOST"),
		NatsPort:                 os.Getenv("NATS_PORT"),
		NatsUser:                 os.Getenv("NATS_USER"),
		NatsPass:                 os.Getenv("NATS_PASS"),
		CreateUserCommandSubject: os.Getenv("CREATE_USER_COMMAND_SUBJECT"),
		CreateUserReplySubject:   os.Getenv("CREATE_USER_REPLY_SUBJECT"),
	}
}
