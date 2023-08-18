// Package server provides functionalities to create and manage HTTP servers.
package server

import "net/http"

// Server represents an HTTP server instance.
type Server struct {
	httpServer *http.Server
}

// Run starts the HTTP server on the specified address using the provided handler.
// It listens for incoming requests and serves them using the given handler.
// The method returns an error if the server fails to start or encounters an issue.
func (s *Server) Run(addr string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.httpServer.ListenAndServe()
}
