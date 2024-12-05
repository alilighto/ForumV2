package server

import (
	"fmt"
	"forum/pkg/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(c *config.API, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + c.Port,
		Handler: handler,
	}
	fmt.Printf("Server is running by: http://%v:%v/\n", c.Host, c.Port)
	return s.httpServer.ListenAndServe()
}
