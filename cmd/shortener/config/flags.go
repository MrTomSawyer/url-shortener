package config

import (
	"flag"
	"os"
)

var ServerAddr string
var DefaultAddr string

func InitServerConfig() {
	flag.StringVar(&ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&DefaultAddr, "b", "http://localhost:8080", "default address and port of a shortened URL")
	flag.Parse()

	if envServerAddr := os.Getenv("SERVER_ADDRESS"); envServerAddr != "" {
		ServerAddr = envServerAddr
	}
	if envDefaultAddr := os.Getenv("BASE_URL"); envDefaultAddr != "" {
		DefaultAddr = envDefaultAddr
	}
}
