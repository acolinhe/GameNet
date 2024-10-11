package test

import (
	"fmt"
	"gamenet/internal/pkg/db"
	"gamenet/internal/pkg/wiki"
	"testing"
)

// Test the entire game pipeline: NER and inserting into the database
func TestGamePipeline(t *testing.T) {
	// Initialize the PostgreSQL connection
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Sample Wikipedia-like text
	text := "The Legend of Zelda is an action-adventure game developed by Nintendo."

	// Run NER on the sample text
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER: %v", err)
	}

	// Test inserting the game with extracted entities into the database
	err = wiki.InsertGameWithEntities(conn, "The Legend of Zelda", text, "1986", entities)
	if err != nil {
		t.Fatalf("Failed to insert game with entities: %v", err)
	}

	// Verify if the game is inserted correctly
	err = verifyGameInsertion(conn, "The Legend of Zelda", entities)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}

	t.Log("Successfully tested the game pipeline, including NER and database insertion.")
}

// Test pipeline with multiple games
func TestGamePipeline_MultipleGames(t *testing.T) {
	// Initialize PostgreSQL connection
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	games := []struct {
		Title       string
		Description string
		ReleaseDate string
	}{
		{"The Legend of Zelda", "The Legend of Zelda is an action-adventure game developed by Nintendo.", "1986"},
		{"Super Mario Bros.", "Super Mario Bros. is a platform game developed by Nintendo.", "1985"},
	}

	for _, game := range games {
		// Run NER on each game's description
		entities, err := wiki.RunNER(game.Description)
		if err != nil {
			t.Fatalf("Failed to run NER for game %s: %v", game.Title, err)
		}

		// Insert each game with its entities into the database
		err = wiki.InsertGameWithEntities(conn, game.Title, game.Description, game.ReleaseDate, entities)
		if err != nil {
			t.Fatalf("Failed to insert game %s with entities: %v", game.Title, err)
		}

		// Verify if the game is inserted correctly
		err = verifyGameInsertion(conn, game.Title, entities)
		if err != nil {
			t.Fatalf("Verification failed for game %s: %v", game.Title, err)
		}
	}

	t.Log("Successfully tested the game pipeline with multiple games.")
}

// Test pipeline with invalid game data
func TestGamePipeline_InvalidData(t *testing.T) {
	// Initialize PostgreSQL connection
	conn, err := db.InitPostgres()
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	// Test with an empty game title
	invalidText := "This is a test description for an invalid game with no title."
	entities, err := wiki.RunNER(invalidText)
	if err != nil {
		t.Fatalf("Failed to run NER: %v", err)
	}

	// Attempt to insert with an invalid game title (empty string)
	err = wiki.InsertGameWithEntities(conn, "", invalidText, "2024", entities)
	if err == nil {
		t.Fatal("Expected failure when inserting game with an empty title, but insertion succeeded.")
	}

	t.Log("Successfully tested error handling for invalid game data.")
}

// Helper function to verify that a game and its entities were correctly inserted into the database
func verifyGameInsertion(conn *db.DB, gameTitle string, expectedEntities []wiki.Entity) error {
	// Verify that the game exists in the database
	var gameID int
	query := `SELECT id FROM Games WHERE title = $1`
	err := conn.QueryRow(query, gameTitle).Scan(&gameID)
	if err != nil {
		return fmt.Errorf("failed to verify game '%s': %v", gameTitle, err)
	}

	// Verify that the related entities were correctly inserted (e.g., developers, platforms, genres)
	for _, entity := range expectedEntities {
		switch entity.Label {
		case "Developer":
			var developerID int
			err = conn.QueryRow(`SELECT id FROM Developers WHERE name = $1`, entity.Text).Scan(&developerID)
			if err != nil {
				return fmt.Errorf("failed to verify developer entity '%s' for game '%s': %v", entity.Text, gameTitle, err)
			}

			// Verify the relationship in the GameDevelopers table
			err = conn.QueryRow(`SELECT 1 FROM GameDevelopers WHERE game_id = $1 AND developer_id = $2`, gameID, developerID).Scan(new(int))
			if err != nil {
				return fmt.Errorf("failed to verify developer relationship for game '%s' and developer '%s': %v", gameTitle, entity.Text, err)
			}
		case "Platform":
			// Similar checks can be done for platforms and other entity types
			var platformID int
			err = conn.QueryRow(`SELECT id FROM Platforms WHERE name = $1`, entity.Text).Scan(&platformID)
			if err != nil {
				return fmt.Errorf("failed to verify platform entity '%s' for game '%s': %v", entity.Text, gameTitle, err)
			}

			// Verify the relationship in the GamePlatforms table
			err = conn.QueryRow(`SELECT 1 FROM GamePlatforms WHERE game_id = $1 AND platform_id = $2`, gameID, platformID).Scan(new(int))
			if err != nil {
				return fmt.Errorf("failed to verify platform relationship for game '%s' and platform '%s': %v", gameTitle, entity.Text, err)
			}
		}
	}

	return nil
}
