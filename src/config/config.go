package config

import (
	"fmt"
	"log"
	"os"
)

const (
	vicinityAgentUrl  = "http://localhost:9997"
	vicinityAdapterID = "967fbc90-c1fa-4390-a438-09a99d2c19cb"
	vicinityVASOid = "7cd7a012-9758-4498-a5c3-bcdbe0ba5c7b"

	serverPort = "9090"

	databasePort = "5432"
	databaseHost = "localhost"
)

type VicinityConfig struct {
	AgentUrl  string
	AdapterID string
	Oid       string
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Host string
	Port string
	User string
	Name string
	Pass string
}

type Config struct {
	Vicinity *VicinityConfig
	Server   *ServerConfig
	Database *DBConfig
}

func (dbc *DBConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbc.Host, dbc.Port, dbc.User, dbc.Name, dbc.Pass,
	)
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Vicinity: &VicinityConfig{
			AgentUrl:  getEnv("VICINITY_AGENT_URL", vicinityAgentUrl),
			AdapterID: getEnv("VICINITY_ADAPTER_ID", vicinityAdapterID),
			Oid: getEnv("VICINITY_VAS_OID", vicinityVASOid),
		},
		Server: &ServerConfig{
			Port: getEnv("SERVER_PORT", serverPort),
		},
		Database: &DBConfig{
			Host: getEnv("DB_HOST", databaseHost),
			Port: getEnv("DB_PORT", databasePort),
			User: getEnv("DB_USER", ""),
			Name: getEnv("DB_NAME", ""),
			Pass: getEnv("DB_PASS", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if isEmpty(defaultVal) {
		log.Printf("environment variable %v is empty\n", key)
		os.Exit(0)
	}

	return defaultVal
}

func isEmpty(val string) bool {
	return val == ""
}
