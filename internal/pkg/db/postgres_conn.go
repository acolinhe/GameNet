package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"os"
)

var DB *sql.DB

// INIT initializes connection to postgres
func INIT() error {
	// Environmental variables
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASS")
	postgresDatabase := os.Getenv("GAMENET_DB")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")

	// Connection string
	conn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword,
		postgresHost, postgresPort, postgresDatabase)

	db, err := sql.Open("postgres", conn)

	// Catch error
	if err != nil {
		return err
	}

	// Ping db to make sure connection works
	if err = db.Ping(); err != nil {
		return err
	}

	fmt.Println("Connected to postgres!")

	return nil
}
