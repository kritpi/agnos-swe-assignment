package main

import (
	"log"

	"github.com/kritpi/agnos-swe-assignment/server"
)

// @title        Agnos SWE Assignment API
// @version      1.0
// @description  REST API for the Agnos SWE assignment (Gin + PostgreSQL).
// @host         localhost:8080
// @BasePath     /api/v1
// @schemes      http
//
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the access token from /staff/login.
func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %+v", err)
	}
}
