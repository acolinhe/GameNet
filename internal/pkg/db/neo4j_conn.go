package db

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
	"os"
)

// Global Neo4j driver for managing connections
var Neo4jDriver neo4j.Driver

// InitNeo4j initializes a connection to the Neo4j database using environment variables.
func InitNeo4j() error {
	// Fetch environment variables for Neo4j connection
	neo4jHost := os.Getenv("NEO4J_HOST")
	neo4jPort := os.Getenv("NEO4J_PORT")
	neo4jUser := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASS")

	// Construct the Neo4j connection URI
	neo4jUri := fmt.Sprintf("bolt://%s:%s", neo4jHost, neo4jPort)

	var err error
	// Initialize the Neo4j driver using the constructed URI and authentication details
	Neo4jDriver, err = neo4j.NewDriver(neo4jUri, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
	if err != nil {
		return err
	}

	// Test the connection to Neo4j to verify connectivity
	if err := Neo4jDriver.VerifyConnectivity(); err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	fmt.Println("Connected to Neo4j!")
	return nil
}

// StoreInNeo4j inserts a game record (title and description) into Neo4j.
// It creates a Game node with the title and description properties.
func StoreInNeo4j(session neo4j.Session, title, description string) error {
	// Execute a write transaction to insert a new Game node into Neo4j
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Cypher query to create a Game node with title and description
		query := `CREATE (g:Game {title: $title, description: $description})`

		// Run the query and pass the parameters
		_, err := tx.Run(query, map[string]interface{}{
			"title":       title,
			"description": description,
		})

		// Return error if the transaction fails
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	// Return error if the transaction couldn't be completed
	if err != nil {
		return fmt.Errorf("could not insert game in Neo4j: %v", err)
	}

	return nil
}

// CloseNeo4j closes the Neo4j driver connection when it's no longer needed
func CloseNeo4j() error {
	if Neo4jDriver != nil {
		// Close the Neo4j driver to release resources
		return Neo4jDriver.Close()
	}
	return nil
}

// Example usage of the functions
func main() {
	// Initialize Neo4j connection
	err := InitNeo4j()
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j: %v", err)
	}
	defer CloseNeo4j() // Ensure Neo4j driver is closed when the program ends

	// Open a session for performing database transactions
	session := Neo4jDriver.NewSession(neo4j.SessionConfig{})
	defer session.Close() // Close the session after usage

	// Insert a game record
	err = StoreInNeo4j(session, "Example Game", "This is an example description")
	if err != nil {
		log.Fatalf("Failed to store game in Neo4j: %v", err)
	}

	fmt.Println("Game inserted successfully.")
}
