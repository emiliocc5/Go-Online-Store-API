package main

import (
	"github.com/emiliocc5/online-store-api/internal/server"
)

func main() {
	s := server.New()

	s.ConfigureRoutes()

	s.Run()
}
