package test

import (
	"database/sql"
	"gamenet/internal/pkg/db"
	"testing"
)

// Test connection to PostgreSQL
func TestPostgresConnection(t *testing.T) {
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Test a sample query to check the connection
	err = conn.Ping()
	if err != nil {
		t.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	t.Log("Successfully connected and pinged PostgreSQL.")
}

// Test inserting a game into the Games table
func TestInsertGame(t *testing.T) {
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Insert a sample game into the Games table
	query := `INSERT INTO Games (title, summary, release_date) VALUES ($1, $2, $3)`
	_, err = conn.Exec(query, "Test Game", "This is a test game.", "2024")
	if err != nil {
		t.Fatalf("Failed to insert game: %v", err)
	}

	t.Log("Successfully inserted a game into the Games table.")
}

// Test retrieving a game from the Games table
func TestRetrieveGame(t *testing.T) {
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Retrieve the game that was inserted
	var title, summary, releaseDate string
	query := `SELECT title, summary, release_date FROM Games WHERE title = $1`
	err = conn.QueryRow(query, "Test Game").Scan(&title, &summary, &releaseDate)
	if err != nil {
		if err == sql.ErrNoRows {
			t.Fatalf("No rows were returned: %v", err)
		}
		t.Fatalf("Failed to retrieve game: %v", err)
	}

	// Check if the data matches what was inserted
	if title != "Test Game" || summary != "This is a test game." || releaseDate != "2024" {
		t.Fatalf("Retrieved data does not match the inserted game.")
	}

	t.Logf("Successfully retrieved game: %s, %s, %s", title, summary, releaseDate)
}

// Test deleting a game from the Games table
func TestDeleteGame(t *testing.T) {
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Delete the game that was inserted
	query := `DELETE FROM Games WHERE title = $1`
	_, err = conn.Exec(query, "Test Game")
	if err != nil {
		t.Fatalf("Failed to delete game: %v", err)
	}

	t.Log("Successfully deleted the game from the Games table.")
}

// Test connection failure with invalid credentials
func TestInvalidPostgresConnection(t *testing.T) {
	// Simulate invalid credentials
	connStr := "host=localhost port=5432 user=invaliduser password=invalidpassword dbname=invaliddb sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err == nil {
		defer conn.Close()
	}

	// Test the connection
	err = conn.Ping()
	if err == nil {
		t.Fatal("Expected failure to connect to PostgreSQL with invalid credentials, but connection succeeded.")
	}

	t.Logf("Failed to connect to PostgreSQL as expected: %v", err)
}

// Test database query failure (e.g., inserting duplicate key)
func TestQueryErrorHandling(t *testing.T) {
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Insert a game to test duplicate insertion
	query := `INSERT INTO Games (title, summary, release_date) VALUES ($1, $2, $3)`
	_, err = conn.Exec(query, "Duplicate Test Game", "This is a test game.", "2024")
	if err != nil {
		t.Fatalf("Failed to insert game: %v", err)
	}

	// Try inserting the same game again to trigger a duplicate key error
	_, err = conn.Exec(query, "Duplicate Test Game", "This is a test game.", "2024")
	if err == nil {
		t.Fatal("Expected failure when inserting duplicate game, but insertion succeeded.")
	}

	t.Logf("Successfully handled duplicate insertion error: %v", err)
}
