package main

import (
	"GameNet/internal/pkg/db"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Test PostgreSQL connection
	err = db.INIT()
	if err != nil {
		log.Fatalf("Error connecting to db: %v", err)
	}

	// fmt.Println("Fetching Wikipedia article...")
	// wiki.FetchWikipediaArticle("Destiny_2")
}
