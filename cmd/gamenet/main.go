package main

import (
	"database/sql"
	"fmt"
	"gamenet/internal/pkg/db"
	"gamenet/internal/pkg/wiki"
	"log"
	"sync"
)

// GameData structure to hold game title, description, and extracted entities
type GameData struct {
	Title       string
	Description string
	Entities    []wiki.Entity
}

func main() {
	// Initialize a connection to the PostgreSQL database
	pgConn, err := db.InitPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgConn.Close() // Close the database connection when the program ends

	// Channels for coordinating between goroutines
	wikiChannel := make(chan wiki.WikiResponse) // Channel to pass fetched Wikipedia data
	nerChannel := make(chan GameData)           // Channel to pass processed NER (Named Entity Recognition) data
	doneChannel := make(chan bool)              // Channel to signal when all tasks are completed

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish

	// Start the goroutine for fetching Wikipedia data
	wg.Add(1)
	go fetchWikipediaData(wikiChannel, &wg) // Fetch Wikipedia data and send it to the wikiChannel

	// Start the goroutine for processing NER on the fetched Wikipedia data
	wg.Add(1)
	go processNER(wikiChannel, nerChannel, &wg) // Process the NER and pass results to nerChannel

	// Start the goroutine for inserting game data and entities into the database
	wg.Add(1)
	go insertGameData(pgConn, nerChannel, doneChannel, &wg) // Insert data into the database, signaled by nerChannel

	// Wait for all goroutines to complete
	wg.Wait()
	close(doneChannel) // Close the doneChannel when all goroutines are finished
	fmt.Println("All tasks completed.")
}

// fetchWikipediaData fetches data from Wikipedia and sends it through the wikiChannel.
// This is executed as a goroutine.
func fetchWikipediaData(wikiChannel chan<- wiki.WikiResponse, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when function completes

	// Fetch data from Wikipedia
	wikiData, err := wiki.FetchWikiData("video_game")
	if err != nil {
		log.Fatalf("Failed to fetch data from Wikipedia: %v", err)
	}

	// Send the fetched data through the channel
	wikiChannel <- *wikiData
	close(wikiChannel) // Close the channel after sending all data
}

// processNER reads data from wikiChannel, processes it for NER, and sends it to nerChannel.
// This is executed as a goroutine.
func processNER(wikiChannel <-chan wiki.WikiResponse, nerChannel chan<- GameData, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when function completes

	// Process each Wikipedia page data from the wikiChannel
	for wikiData := range wikiChannel {
		for _, page := range wikiData.Query.Pages {
			title := page.Title
			description := page.Extract

			// Run NER (Named Entity Recognition) on the page description
			entities, err := wiki.RunNER(description)
			if err != nil {
				log.Printf("Failed to run NER on %s: %v", title, err)
				continue
			}

			// Send the processed game data (with entities) to the nerChannel
			nerChannel <- GameData{
				Title:       title,
				Description: description,
				Entities:    entities,
			}
		}
	}
	// Close the nerChannel after all data has been processed
	close(nerChannel)
}

// InsertGameWithEntities inserts the game and its related entities (Developers, Platforms, Genres) into the database concurrently.
func InsertGameWithEntities(db *sql.DB, title, summary, releaseDate string, entities []Entity) error {
	// Insert the game into the database and get the gameID
	gameID, err := insertGame(db, title, summary, releaseDate)
	if err != nil {
		return fmt.Errorf("failed to insert game: %v", err)
	}

	var wg sync.WaitGroup                      // WaitGroup to track goroutines inserting entities
	errChan := make(chan error, len(entities)) // Channel to collect any errors from the goroutines

	// For each entity (e.g., Developer, Platform, Genre), insert it concurrently
	for _, entity := range entities {
		wg.Add(1) // Increment WaitGroup counter for each entity

		go func(entity Entity) {
			defer wg.Done() // Mark this goroutine as done after the entity is processed

			var insertErr error
			// Insert the entity based on its type (Developer, Platform, Genre)
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

			// If there is an error, send it to the error channel
			if insertErr != nil {
				errChan <- fmt.Errorf("failed to insert entity (%s): %v", entity.Text, insertErr)
			}
		}(entity) // Pass the current entity to the goroutine
	}

	wg.Wait()      // Wait for all entity-inserting goroutines to finish
	close(errChan) // Close the error channel after all goroutines have finished

	// Collect any errors from the error channel
	var errorMessages []string
	for err := range errChan {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}

	// If there were any errors, return them
	if len(errorMessages) > 0 {
		return fmt.Errorf("multiple errors occurred: %v", errorMessages)
	}

	return nil
}
