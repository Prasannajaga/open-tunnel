package config

import "opentunnel/server/constants"

type Config struct {
	ControlPort  string
	ExternalPort string
	DataPort     string
	HTTPPort     string
	JWTSecret    string
	TokenExpiry  int
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
}

func NewConfig() *Config {
	return &Config{
		ControlPort:  constants.ControlPort,
		ExternalPort: constants.ExternalPort,
		DataPort:     constants.DataPort,
		HTTPPort:     constants.HTTPPort,
		JWTSecret:    "opentunnel-secret-key-change-in-production",
		TokenExpiry:  24,
		DBHost:       "localhost",
		DBPort:       "5432",
		DBUser:       "postgres",
		DBPassword:   "postgres",
		DBName:       "opentunnel",
	}
}
