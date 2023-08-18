// Package config provides structures and functions for managing application configuration.
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

// AppConfig holds the configuration parameters for the application.
type AppConfig struct {
	Server struct {
		DefaultAddr string `env:"BASE_URL"`     // DefaultAddr is the base URL for the application.
		ServerAddr  string `env:"SERVER_ADDRESS"` // ServerAddr is the address and port to run the server.
		TempFolder  string `env:"FILE_STORAGE_PATH"` // TempFolder is the path for temporary file storage.
		SecretKey   string `env:"SECRET_KEY"`      // SecretKey is the secret signing key.
	}
	DataBase struct {
		ConnectionStr string `env:"DATABASE_DSN"` // ConnectionStr is the database connection string.
	}
}

// InitAppConfig initializes the application configuration by parsing command-line flags and environment variables.
func (a *AppConfig) InitAppConfig() error {
	flag.StringVar(&a.Server.ServerAddr, "a", ":8080", "address and port to run the server")
	flag.StringVar(&a.Server.DefaultAddr, "b", "http://localhost:8080", "default address and port of a shortened URL")
	flag.StringVar(&a.Server.TempFolder, "f", "/tmp/short-url-db.json", "default temp data storage path and filename")
	flag.StringVar(&a.Server.SecretKey, "sk", "secret", "secret signing key")
	flag.StringVar(&a.DataBase.ConnectionStr, "d", "", "Database connection string")

	// Uncomment the following line for local testing
	// flag.StringVar(&a.DataBase.ConnectionStr, "d", "host=localhost port=5432 user=myuser password=password dbname=mydb sslmode=disable", "Database connection string")

	flag.Parse()

	err := env.Parse(a)
	if err != nil {
		return err
	}
	return nil
}
