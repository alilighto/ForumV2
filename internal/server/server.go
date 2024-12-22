package server

import (
	"log"
	"net/http"

	"forum/pkg/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(c *config.API, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    c.Host + ":" + c.Port,
		Handler: handler,
	}
	log.Printf("\033[32m Server is running...🚀\nLink: 🌐 http://%s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
