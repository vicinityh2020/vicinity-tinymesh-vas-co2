package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	vicinityAgentUrl  = "http://localhost:9997"
	vicinityAdapterID = "967fbc90-c1fa-4390-a438-09a99d2c19cb"
	vicinityVASOid    = "7cd7a012-9758-4498-a5c3-bcdbe0ba5c7b"
	vicinityKPIKey    = ""

	serverPort = "9090"

	databasePort = "5432"
	databaseHost = "localhost"

	smsSender    = "CWi:Moss CO2"

	// :-)
	smsRecipient = "41369367"
)

type VicinityConfig struct {
	AgentUrl  string
	AdapterID string
	Oid       string
	KPIKey    string
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

type SMSConfig struct {
	User       string
	Key        string
	Sender     string
	Recipients []string
}

type Config struct {
	Vicinity *VicinityConfig
	Server   *ServerConfig
	Database *DBConfig
	SMS      *SMSConfig
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
			Oid:       getEnv("VICINITY_VAS_OID", vicinityVASOid),
			KPIKey:    getEnv("VICINITY_KPI_KEY", vicinityKPIKey),
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
		SMS: &SMSConfig{
			Recipients: strings.Split(getEnv("KEYSMS_RECIPIENTS", smsRecipient), ","),
			Sender:     getEnv("KEYSMS_SENDER", smsSender),
			User:       getEnv("KEYSMS_USER", ""),
			Key:        getEnv("KEYSMS_API_KEY", ""),
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
