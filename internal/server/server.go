package server

import (
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/emiliocc5/online-store-api/package/handler"
	"github.com/gin-gonic/gin"
)

//TODO Implement Singleton

type Server struct {
	server *gin.Engine
}

func New() Server {
	s := &Server{}
	g := gin.Default()

	s.server = g

	//TODO Take this to main
	ch := handler.CartHandlerImpl{
		CartService: &services.CartServiceImpl{},
	}
	s.server.GET("/api/cart", ch.HandleGetCart)

	return *s
}

func (s *Server) Run() {
	err := s.server.Run()
	if err != nil {
		return
	}
}
