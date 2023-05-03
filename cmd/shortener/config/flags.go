package config

import "flag"

var ServerAddr string
var DefaultAddr string

func ParseFlags() {
	flag.StringVar(&ServerAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.StringVar(&DefaultAddr, "b", "http://localhost:8080", "default address and port of a shortened URL")
	flag.Parse()
}
