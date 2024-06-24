package main

import (
	"GameNet/internal/pkg/db"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
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

	// Test Neo4j connection
	err = db.InitNeo4j()
	if err != nil {
		log.Fatalf("Error connecting to Neo4j: %v", err)
	}
	defer func(Neo4jDriver neo4j.Driver) {
		err := Neo4jDriver.Close()
		if err != nil {

		}
	}(db.Neo4jDriver)

	// fmt.Println("Fetching Wikipedia article...")
	// wiki.FetchWikipediaArticle("Destiny_2")
}
