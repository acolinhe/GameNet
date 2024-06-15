package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

const wikipediaAPIURL = "https://en.wikipedia.org/w/api.php"

type WikipediaResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]struct {
			Pageid  int    `json:"pageid"`
			Ns      int    `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wikipediaToken := os.Getenv("WIKIPEDIA_TOKEN")
	if wikipediaToken == "" {
		log.Fatal("Wikipedia token not set")
	}

	testArticle := "Destiny_2"
	url := fmt.Sprintf("%s?action=query&format=json&prop=extracts&exintro&titles=%s", wikipediaAPIURL, testArticle)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", wikipediaToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var wikiResponse WikipediaResponse
	if err := json.Unmarshal(body, &wikiResponse); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	for _, page := range wikiResponse.Query.Pages {
		fmt.Printf("Title: %s\n\nExtract:\n%s\n\n", page.Title, page.Extract)
	}
}
