package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Username string
	Port     string
	DBName   string
	SSLMode  string

	Password    string
	MaxOpenConn int8
	MaxIdleConn int8

	// Server configuration
	ServerPort string
}

func (s *Config) LoadEnv() (*Config, error) {
	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // default port
	}

	return &Config{
		Host:        os.Getenv("HOST"),
		Username:    os.Getenv("USERNAME"),
		Port:        os.Getenv("PORT"),
		DBName:      os.Getenv("DB_NAME"),
		Password:    password,
		SSLMode:     os.Getenv("SSL_MODE"),
		MaxOpenConn: 8,
		MaxIdleConn: 5,
		ServerPort:  serverPort,
	}, nil
}

func (s *Config) BuildDSN(config *Config) (string, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Host, config.Username, config.Password, config.DBName, config.Port, config.SSLMode)
	return dsn, nil
}
