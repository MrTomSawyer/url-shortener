package main

import (
	"net/http"

	h "github.com/MrTomSawyer/url-shortener/cmd/shortener/httphandlers"
)

func run(m *http.ServeMux) error {
	return http.ListenAndServe(`:8080`, m)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HTTPHandler)

	if err := run(mux); err != nil {
		panic(err)
	}
}
