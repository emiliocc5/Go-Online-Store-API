package main

import (
	"github.com/emiliocc5/online-store-api/internal/server"
)

func main() {
	s := server.New()

	/*r := s.Router()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})*/

	s.Run()
}
