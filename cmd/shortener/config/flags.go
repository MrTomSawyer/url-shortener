package config

import "flag"

var ServerAddr string
var DefaultAddr string

func ParseFlags() {
	flag.StringVar(&ServerAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&DefaultAddr, "b", "localhost:8080", "address and port to run server")
	flag.Parse()
}
