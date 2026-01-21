package config

import (
	"os"
	"strconv"

	"opentunnel/client/constants"
)

type Config struct {
	ServerIP     string
	ControlPort  string
	ExternalPort string
	DataPort     string
	HTTPPort     string
	LocalPort    int
}

func NewConfig() *Config {
	return &Config{
		ServerIP:     getEnv("SERVER_IP", "34.133.55.212"),
		ControlPort:  getEnv("CONTROL_PORT", constants.ControlPort),
		ExternalPort: getEnv("EXTERNAL_PORT", constants.ExternalPort),
		DataPort:     getEnv("DATA_PORT", constants.DataPort),
		HTTPPort:     getEnv("HTTP_PORT", "8081"),
		LocalPort:    getEnvInt("LOCAL_PORT", constants.DefaultLocalPort),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func (c *Config) ServerAddr() string {
	return c.ServerIP + ":" + c.ControlPort
}

func (c *Config) DataAddr() string {
	return c.ServerIP + ":" + c.DataPort
}

func (c *Config) ExposedURL() string {
	return "http://" + c.ServerIP + ":" + c.ExternalPort
}
