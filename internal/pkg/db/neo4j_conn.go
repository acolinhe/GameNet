package db

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"os"
)

var Neo4jDriver neo4j.Driver

func InitNeo4j() error {
	neo4jHost := os.Getenv("NEO4J_HOST")
	neo4jPort := os.Getenv("NEO4J_PORT")
	neo4jUser := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASS")

	neo4jUri := fmt.Sprintf("bolt://%s:%s", neo4jHost, neo4jPort)

	var err error
	Neo4jDriver, err = neo4j.NewDriver(neo4jUri, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""))
	if err != nil {
		return err
	}

	// Test the connection
	if err := Neo4jDriver.VerifyConnectivity(); err != nil {
		return fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	fmt.Println("Connected to Neo4j!")
	return nil
}
