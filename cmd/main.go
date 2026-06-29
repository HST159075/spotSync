package main

import (
	"log"
	"os"
	"spotsync/internal/config"
	"spotsync/internal/server"
	"github.com/joho/godotenv"
)
func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config.ConnectDatabase()

	e := server.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}