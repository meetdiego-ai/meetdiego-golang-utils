package fetcher

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// DefaultHeaders provides common headers for HTTP requests
var DefaultHeaders = map[string]string{
	"User-Agent": "Mozilla/5.0 (iPad; CPU OS 12_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148",
}

func NormalizeURL(rawURL string) (string, error) {
	// Trim spaces first
	rawURL = strings.TrimSpace(rawURL)

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Convert to lowercase
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	// Remove www. prefix
	parsedURL.Host = strings.TrimPrefix(parsedURL.Host, "www.")

	// Remove query parameters and fragments
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""

	// Normalize the path
	if parsedURL.Path == "" || parsedURL.Path == "/" {
		parsedURL.Path = "/" // Ensure root URL is represented as "/"
	} else {
		// Remove trailing slashes
		parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	}

	// Clean up repeated slashes
	parsedURL.Path = strings.Join(strings.Fields(parsedURL.Path), "/")

	return parsedURL.String(), nil
}

func IsValidHtmlPage(pageURL string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", pageURL, nil)
	if err != nil {
		return false, fmt.Errorf("invalid URL (NewRequest): %v", err)
	}

	// Apply default headers
	for key, value := range DefaultHeaders {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("invalid URL (Do): %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		return false, fmt.Errorf("URL returned status code %d", resp.StatusCode)
	}

	// Check content type is html
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return false, fmt.Errorf("URL returned content type %s", resp.Header.Get("Content-Type"))
	}

	return true, nil
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

func ReadPageContent(pageURL string) (string, error) {
	resp, err := FetchURL(pageURL, nil)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	// Strip the document of certain tags such as script, style, link, meta, etc.
	doc.Find("script, style, link, meta").Each(func(index int, s *goquery.Selection) {
		s.Remove()
	})

	content := doc.Text()

	return content, nil
}
