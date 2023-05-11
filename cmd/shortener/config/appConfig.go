package config

import (
	"flag"
	"os"
)

type AppConfig struct {
	Server struct {
		DefaultAddr string
		ServerAddr  string
	}
}

func InitAppConfig(config *AppConfig) {
	flag.StringVar(&config.Server.ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.Server.DefaultAddr, "b", "http://localhost:8080", "default address and port of a shortened URL")

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		config.Server.ServerAddr = envServerAddr
	}
	if envDefaultAddr := os.Getenv("BASE_URL"); envDefaultAddr != "" {
		config.Server.DefaultAddr = envDefaultAddr
	}
}
