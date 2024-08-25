package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SpellCheckError represents the spelling error info
type SpellCheckError struct {
	Code int      `json:"code"` // error code
	Pos  int      `json:"pos"`  // position of the error in the text
	Row  int      `json:"row"`  // row number
	Col  int      `json:"col"`  // column number
	Len  int      `json:"len"`  // length of the incorrect word
	Word string   `json:"word"` // the word with error
	S    []string `json:"s"`    // suggested corrections
}

// SpellChecker is an interface for spell checking service
type SpellChecker interface {
	Check(text string) ([]SpellCheckError, error)
}

// YandexSpellChecker is an implementation of SpellChecker interface (logic)
type YandexSpellChecker struct {
	client *http.Client
}

// NewYandexSpellChecker creates a custom HTTP client with timeout
func NewYandexSpellChecker() *YandexSpellChecker {
	return &YandexSpellChecker{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Check verifies the text for spelling errors using the external API
func (ysc *YandexSpellChecker) Check(text string) ([]SpellCheckError, error) {
	apiURL := os.Getenv("SPELLCHECK_API_URL")

	if apiURL == "" {
		return nil, fmt.Errorf("API URL is not set in the environment variables")
	}

	// Prepare the request parameters
	data := url.Values{}
	data.Set("text", text)
	data.Set("lang", "ru,en")

	// New POST request to the external API
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := ysc.client.Do(req) // Send the request to the API
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SpellCheck API returned status %d: %s", resp.StatusCode, body)
	}

	// Read and parse the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var spellErrors []SpellCheckError
	if err := json.Unmarshal(body, &spellErrors); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return spellErrors, nil
}
