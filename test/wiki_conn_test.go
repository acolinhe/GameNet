package test

import (
	"gamenet/internal/pkg/wiki"
	"testing"
)

// Test NER on normal text with multiple entities
func TestRunNER(t *testing.T) {
	// Sample text for testing
	text := "Nintendo is a video game company that developed Mario."

	// Call the RunNER function
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER: %v", err)
	}

	// Check if entities were extracted
	if len(entities) == 0 {
		t.Fatalf("No entities extracted from the text.")
	}

	// Check for the "Nintendo" entity as a Developer
	found := false
	for _, entity := range entities {
		if entity.Label == "Developer" && entity.Text == "Nintendo" {
			found = true
		}
	}
	if !found {
		t.Fatalf("Expected entity 'Nintendo' not found in the text.")
	}

	t.Log("Successfully extracted entities from text.")
}

// Test NER on empty input
func TestRunNER_EmptyInput(t *testing.T) {
	// Test with empty input
	text := ""

	// Call the RunNER function
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER on empty input: %v", err)
	}

	// Check that no entities are returned for empty input
	if len(entities) != 0 {
		t.Fatalf("Expected no entities, but got: %v", entities)
	}

	t.Log("Successfully handled empty input.")
}

// Test NER on input with no recognizable entities
func TestRunNER_NoEntities(t *testing.T) {
	// Input that has no recognizable entities
	text := "This is just a random sentence with no entities."

	// Call the RunNER function
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER on input with no entities: %v", err)
	}

	// Check that no entities are extracted
	if len(entities) != 0 {
		t.Fatalf("Expected no entities, but got: %v", entities)
	}

	t.Log("Successfully handled input with no entities.")
}

// Test NER on input with multiple types of entities
func TestRunNER_MultipleEntities(t *testing.T) {
	// Input with multiple types of entities (Developer, Game, Platform)
	text := "Nintendo developed the game Super Mario on the Nintendo Switch platform."

	// Call the RunNER function
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER: %v", err)
	}

	// Expected entities
	expectedEntities := map[string]string{
		"Nintendo":        "Developer",
		"Super Mario":     "Game",
		"Nintendo Switch": "Platform",
	}

	// Check that each expected entity is extracted correctly
	for expectedText, expectedLabel := range expectedEntities {
		found := false
		for _, entity := range entities {
			if entity.Text == expectedText && entity.Label == expectedLabel {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected entity '%s' with label '%s' not found in the text.", expectedText, expectedLabel)
		}
	}

	t.Log("Successfully extracted multiple entities from text.")
}

// Test NER on input with special characters
func TestRunNER_SpecialCharacters(t *testing.T) {
	// Input with special characters
	text := "Apple Inc. developed the iPhone, which runs on iOS."

	// Call the RunNER function
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER: %v", err)
	}

	// Expected entities
	expectedEntities := map[string]string{
		"Apple Inc.": "Developer",
		"iPhone":     "Product",
		"iOS":        "Platform",
	}

	// Check that each expected entity is extracted correctly
	for expectedText, expectedLabel := range expectedEntities {
		found := false
		for _, entity := range entities {
			if entity.Text == expectedText && entity.Label == expectedLabel {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected entity '%s' with label '%s' not found in the text.", expectedText, expectedLabel)
		}
	}

	t.Log("Successfully extracted entities with special characters.")
}

// Test NER with large input
func TestRunNER_LargeInput(t *testing.T) {
	// Large input text simulating a long article
	text := `
		Google is a multinational company that develops products such as Android, Google Cloud, and Chrome.
		Microsoft, another major player, developed the Windows OS and Azure cloud platform. 
		Apple Inc. is well known for its products such as the iPhone and MacBook.
	`

	// Call the RunNER function
	entities, err := wiki.RunNER(text)
	if err != nil {
		t.Fatalf("Failed to run NER on large input: %v", err)
	}

	// Expected entities
	expectedEntities := map[string]string{
		"Google":       "Developer",
		"Android":      "Product",
		"Google Cloud": "Platform",
		"Microsoft":    "Developer",
		"Windows":      "Product",
		"Azure":        "Platform",
		"Apple Inc.":   "Developer",
		"iPhone":       "Product",
		"MacBook":      "Product",
	}

	// Check that each expected entity is extracted correctly
	for expectedText, expectedLabel := range expectedEntities {
		found := false
		for _, entity := range entities {
			if entity.Text == expectedText && entity.Label == expectedLabel {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected entity '%s' with label '%s' not found in the text.", expectedText, expectedLabel)
		}
	}

	t.Log("Successfully extracted entities from large input.")
}
