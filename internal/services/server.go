package services

import "github.com/gin-gonic/gin"

type Server struct {
	server *gin.Engine
}

func New() Server {
	s := &Server{}
	g := gin.Default()

	s.server = g

	return *s
}

func (s *Server) Router() *gin.Engine {
	return s.server
}

func (s *Server) Run() {
	err := s.server.Run()
	if err != nil {
		return
	}
}
