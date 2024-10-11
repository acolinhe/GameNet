package wiki

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
)

// Entity represents a single recognized entity from NER (Named Entity Recognition)
type Entity struct {
	Text  string `json:"text"`  // The entity text extracted from the input text
	Label string `json:"label"` // The type of entity (e.g., Developer, Platform, Genre)
}

// RunNER executes a Python script to perform Named Entity Recognition on a given text
// and returns a list of recognized entities.
func RunNER(text string) ([]Entity, error) {
	// Command to run the Python NER script with the input text
	cmd := exec.Command("python3", "ner.py", text)

	// Buffer to capture the script's output
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command and check for errors
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run Python NER script: %v", err)
	}

	// Parse the JSON output from the Python script into a slice of Entity
	var entities []Entity
	if err := json.Unmarshal(out.Bytes(), &entities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Python NER output: %v", err)
	}

	// Return the list of recognized entities
	return entities, nil
}

// InsertGameWithEntitiesWithContext inserts a game and its related entities (Developers, Platforms, Genres)
// into the database concurrently with context cancellation support.
func InsertGameWithEntitiesWithContext(ctx context.Context, db *sql.DB, title, summary, releaseDate string, entities []Entity) error {
	// Insert the game into the Games table and retrieve its gameID
	gameID, err := insertGame(db, title, summary, releaseDate)
	if err != nil {
		return fmt.Errorf("failed to insert game: %v", err)
	}

	// WaitGroup to synchronize the concurrent insertion of entities
	var wg sync.WaitGroup
	// Buffered channel to collect errors from the goroutines
	errChan := make(chan error, len(entities))

	// Process each entity concurrently
	for _, entity := range entities {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine

		// Start a goroutine to insert each entity (Developer, Platform, Genre)
		go func(entity Entity) {
			defer wg.Done() // Decrement the WaitGroup counter when this goroutine finishes

			// Check if the context has been canceled before proceeding
			select {
			case <-ctx.Done():
				// Exit the goroutine if context is canceled
				return
			default:
				// Insert entity based on its label
				var insertErr error
				switch entity.Label {
				case "Developer":
					insertErr = insertDeveloper(db, gameID, entity.Text)
				case "Platform":
					insertErr = insertPlatform(db, gameID, entity.Text)
				case "Genre":
					insertErr = insertGenre(db, gameID, entity.Text)
				default:
					return
				}

				// Send any insertion error to the error channel
				if insertErr != nil {
					errChan <- fmt.Errorf("failed to insert entity (%s): %v", entity.Text, insertErr)
				}
			}
		}(entity) // Pass the entity to the anonymous function
	}

	// Wait for all entity-inserting goroutines to complete
	wg.Wait()
	close(errChan) // Close the error channel when all goroutines have finished

	// Collect all errors from the error channel
	var errorMessages []string
	for err := range errChan {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	// If there were any errors, return them as a single error
	if len(errorMessages) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errorMessages)
	}

	// Return nil if no errors occurred
	return nil
}

// insertGame inserts a game into the Games table and returns the generated gameID.
func insertGame(db *sql.DB, title, summary, releaseDate string) (int, error) {
	// SQL query to insert the game and return the generated game ID
	query := `INSERT INTO Games (title, summary, release_date) VALUES ($1, $2, $3) RETURNING id`
	var gameID int
	// Execute the query and scan the generated ID into gameID
	err := db.QueryRow(query, title, summary, releaseDate).Scan(&gameID)
	if err != nil {
		return 0, err
	}
	// Return the generated gameID
	return gameID, nil
}

// insertDeveloper inserts a developer into the Developers table (if it doesn't exist)
// and adds a record to the GameDevelopers table for the relationship.
func insertDeveloper(db *sql.DB, gameID int, developerName string) error {
	// Query to check if the developer already exists in the Developers table
	var developerID int
	query := `SELECT id FROM Developers WHERE name = $1`
	err := db.QueryRow(query, developerName).Scan(&developerID)

	if err == sql.ErrNoRows {
		// If the developer doesn't exist, insert it into the Developers table
		err = db.QueryRow(`INSERT INTO Developers (name) VALUES ($1) RETURNING id`, developerName).Scan(&developerID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Insert the relationship between the game and the developer into the GameDevelopers table
	_, err = db.Exec(`INSERT INTO GameDevelopers (game_id, developer_id) VALUES ($1, $2)`, gameID, developerID)
	return err
}

// insertPlatform inserts a platform into the Platforms table (if it doesn't exist)
// and adds a record to the GamePlatforms table for the relationship.
func insertPlatform(db *sql.DB, gameID int, platformName string) error {
	// Query to check if the platform already exists in the Platforms table
	var platformID int
	query := `SELECT id FROM Platforms WHERE name = $1`
	err := db.QueryRow(query, platformName).Scan(&platformID)

	if err == sql.ErrNoRows {
		// If the platform doesn't exist, insert it into the Platforms table
		err = db.QueryRow(`INSERT INTO Platforms (name) VALUES ($1) RETURNING id`, platformName).Scan(&platformID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Insert the relationship between the game and the platform into the GamePlatforms table
	_, err = db.Exec(`INSERT INTO GamePlatforms (game_id, platform_id) VALUES ($1, $2)`, gameID, platformID)
	return err
}

// insertGenre inserts a genre into the Genres table (if it doesn't exist)
// and adds a record to the GameGenres table for the relationship.
func insertGenre(db *sql.DB, gameID int, genreName string) error {
	// Query to check if the genre already exists in the Genres table
	var genreID int
	query := `SELECT id FROM Genres WHERE name = $1`
	err := db.QueryRow(query, genreName).Scan(&genreID)

	if err == sql.ErrNoRows {
		// If the genre doesn't exist, insert it into the Genres table
		err = db.QueryRow(`INSERT INTO Genres (name) VALUES ($1) RETURNING id`, genreName).Scan(&genreID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Insert the relationship between the game and the genre into the GameGenres table
	_, err = db.Exec(`INSERT INTO GameGenres (game_id, genre_id) VALUES ($1, $2)`, gameID, genreID)
	return err
}
