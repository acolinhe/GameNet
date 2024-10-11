package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
	"os"
	"time"
)

// PostgreSQL connection details - update these as per your environment.
const (
	host     = "localhost"
	port     = 5432
	user     = "yourusername"
	password = "yourpassword"
	dbname   = "gamenet"
)

// InitPostgres initializes a connection to the PostgreSQL database and configures connection pooling.
// It returns the *sql.DB object representing the connection and an error if any.
func InitPostgres() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_DB_HOST")
	user := os.Getenv("POSTGRES_DB_USER")
	password := os.Getenv("POSTGRES_DB_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_NAME")
	port := os.Getenv("POSTGRES_DB_PORT")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set connection pooling options
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}

// StoreInPostgres inserts a game into the 'games' table in the PostgreSQL database.
// The function takes a connection, the title, and description of the game.
func StoreInPostgres(db *sql.DB, title, description string) error {
	// SQL query to insert the game into the 'games' table
	query := `INSERT INTO games (title, description) VALUES ($1, $2)`

	// Execute the insert query, passing the title and description as parameters
	_, err := db.Exec(query, title, description)
	if err != nil {
		return fmt.Errorf("could not insert game: %v", err)
	}
	return nil
}

// ClosePostgres closes the connection to PostgreSQL when it's no longer needed.
// This helps in cleaning up resources.
func ClosePostgres(db *sql.DB) error {
	if db != nil {
		// Close the database connection
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close PostgreSQL connection: %v", err)
		}
	}
	return nil
}
