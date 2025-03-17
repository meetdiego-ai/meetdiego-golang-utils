package serpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// Define a struct to match the JSON format in test.json
// ... existing code ...

type SearchResult struct {
	SearchMetadata   SearchMetadata    `json:"search_metadata"`
	RelatedQuestions []RelatedQuestion `json:"related_questions"`
	OrganicResults   []OrganicResult   `json:"organic_results"`
	AnswerBox        AnswerBox         `json:"answer_box"`
}

type RelatedQuestion struct {
	Question string `json:"question"`
	Snippet  string `json:"snippet"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Link     string `json:"link"`
	UUID     string `json:"uuid"`
}

type OrganicResult struct {
	Position int    `json:"position"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Snippet  string `json:"snippet"`
	Source   string `json:"source"`
	UUID     string `json:"uuid"`
}

type AnswerBox struct {
	Title   string   `json:"title"`
	Link    string   `json:"link"`
	Snippet string   `json:"snippet"`
	Source  string   `json:"source"`
	Images  []string `json:"images"`
	UUID    string   `json:"uuid"`
}

// Structs to match the JSON format in test.json
// ... existing code ...
type SearchMetadata struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	JsonEndpoint   string  `json:"json_endpoint"`
	CreatedAt      string  `json:"created_at"`
	ProcessedAt    string  `json:"processed_at"`
	GoogleURL      string  `json:"google_url"`
	RawHTMLFile    string  `json:"raw_html_file"`
	TotalTimeTaken float64 `json:"total_time_taken"`
}

// ... existing code ...

// Add more structs as needed to match the JSON structure

// generateUUIDs adds UUIDs to all items in the search result
func generateUUIDs(result *SearchResult) {
	for i := range result.OrganicResults {
		result.OrganicResults[i].UUID = uuid.New().String()
	}

	for i := range result.RelatedQuestions {
		result.RelatedQuestions[i].UUID = uuid.New().String()
	}

	result.AnswerBox.UUID = uuid.New().String()
}

// KeywordSearch calls the SerpApi and returns the search results
func KeywordSearch(keyword string) (SearchResult, error) {
	apiKey := os.Getenv("SERPAPI_KEY")
	if apiKey == "" {
		fmt.Println("API key not found")
		return SearchResult{}, fmt.Errorf("API key not found")
	}

	// Need to short circuit here if there is a file in the current directory called test.json
	// and return the contents of that file
	if _, err := os.Stat("test.json"); err == nil {
		body, err := os.ReadFile("test.json")
		if err != nil {
			return SearchResult{}, err
		}

		var result SearchResult
		err = json.Unmarshal(body, &result)
		if err != nil {
			return SearchResult{}, err
		}

		generateUUIDs(&result)
		return result, nil
	}

	fmt.Println("Calling SerpAPI with keyword:", keyword)
	body, err := callSerpApi(keyword, apiKey)
	if err != nil {
		return SearchResult{}, err
	}

	var result SearchResult

	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return SearchResult{}, err
	}

	generateUUIDs(&result)
	return result, nil
}

// callSerpApi makes the HTTP request to SerpAPI and returns the response body
func callSerpApi(keyword string, apiKey string) ([]byte, error) {
	url := fmt.Sprintf("https://serpapi.com/search.json?q=%s&api_key=%s", keyword, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting response:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	return body, nil
}
