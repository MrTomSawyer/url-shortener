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

func (a *AppConfig) InitAppConfig() {
	flag.StringVar(&a.Server.ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&a.Server.DefaultAddr, "b", "http://localhost:8080", "default address and port of a shortened URL")

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		a.Server.ServerAddr = envServerAddr
	}
	if envDefaultAddr := os.Getenv("BASE_URL"); envDefaultAddr != "" {
		a.Server.DefaultAddr = envDefaultAddr
	}
}
