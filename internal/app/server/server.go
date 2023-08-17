package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(addr string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return s.httpServer.ListenAndServe()
}
