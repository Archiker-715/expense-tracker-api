package main

import (
	"log"

	"github.com/joho/godotenv"
)

// TODO: add routing and middleware for service's apipaths
// init DB
// dockerize that all
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment variables")
	}
}
