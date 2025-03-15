package fetcher

import (
	"fmt"
	"io"
	"net/http"
)

// DefaultHeaders provides common headers for HTTP requests
var DefaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
}

// FetchURL retrieves content from a URL with custom headers
func FetchURL(pageURL string, customHeaders map[string]string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}

	// Apply default headers
	for key, value := range DefaultHeaders {
		req.Header.Set(key, value)
	}

	// Apply custom headers (overriding defaults if needed)
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 304 {
		return nil, fmt.Errorf("URL returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
