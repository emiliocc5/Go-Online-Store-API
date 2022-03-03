package main

import (
	"github.com/emiliocc5/online-store-api/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	s := services.New()

	r := s.Router()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	s.Run()
}
