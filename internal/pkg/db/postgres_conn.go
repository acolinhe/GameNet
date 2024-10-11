package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
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
	// Connection string with credentials for PostgreSQL connection
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a new database connection using the PostgreSQL driver
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %v", err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(10)                  // Maximum number of open connections
	db.SetMaxIdleConns(5)                   // Maximum number of idle connections
	db.SetConnMaxLifetime(30 * time.Minute) // Maximum connection lifetime

	// Test the connection to make sure it's active
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
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
