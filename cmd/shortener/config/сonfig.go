package config

import (
	"flag"
	"os"
)

type AppConfig struct {
	Server struct {
		DefaultAddr string
		ServerAddr  string
		TempFolder  string
	}
}

func (a *AppConfig) InitAppConfig() {
	flag.StringVar(&a.Server.ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&a.Server.DefaultAddr, "b", "http://localhost:8080", "default address and port of a shortened URL")
	flag.StringVar(&a.Server.TempFolder, "f", "/tmp/short-url-db.json", "default temp data storage path and filename")

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		a.Server.ServerAddr = envServerAddr
	}
	if envDefaultAddr := os.Getenv("BASE_URL"); envDefaultAddr != "" {
		a.Server.DefaultAddr = envDefaultAddr
	}
	if envTempFolder := os.Getenv("FILE_STORAGE_PATH"); envTempFolder != "" {
		a.Server.TempFolder = envTempFolder
	}
}
