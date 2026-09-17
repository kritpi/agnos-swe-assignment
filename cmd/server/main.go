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
func main() {
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %+v", err)
	}
}
