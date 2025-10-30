package gopastemyst

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Check https://docs.beta.myst.rs/pastes

// TODO: Add documentation to each struct type and function

type APIError struct {
	StatusMessage string `json:"statusMessage"`
}

func (e APIError) Error() string {
	return e.StatusMessage
}

type Pasty struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

type Paste struct {
	// Non-Nullable fields
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresIn string    `json:"expiresIn"`
	Pinned    bool      `json:"pinned"`
	Private   bool      `json:"private"`
	Stars     int       `json:"stars"`
	Tags      []string  `json:"tags"`
	Pasties   []Pasty   `json:"pasties"`

	// Nullable fields
	OwnerID   *string    `json:"ownerID"`
	EditedAt  *time.Time `json:"editedAt"`
	DeletesAt *time.Time `json:"deletesAt"`
}

type PastyStats struct {
	Bytes int `json:"bytes"`
	Lines int `json:"lines"`
	Words int `json:"words"`
}

type Stats struct {
	Bytes   int                   `json:"bytes"`
	Lines   int                   `json:"lines"`
	Pasties map[string]PastyStats `json:"pasties"`
	Words   int                   `json:"words"`
}

type CreatePastyOptions struct {
	Title    string `json:"title,omitempty"`
	Content  string `json:"content"`
	Language string `json:"language,omitempty"`
}

type CreatePasteOptions struct {
	Title     string               `json:"title,omitempty"`
	ExpiresIn string               `json:"expiresIn,omitempty"`
	Anonymous bool                 `json:"anonymous,omitempty"`
	Private   bool                 `json:"private,omitempty"`
	Pinned    bool                 `json:"pinned,omitempty"`
	Encrypted bool                 `json:"encrypted,omitempty"`
	Tags      []string             `json:"tags,omitempty"`
	Pasties   []CreatePastyOptions `json:"pasties"`
}

type LanguageStats struct {
	Aliases            []string `json:"aliases"`
	CodemirrorMimeType string   `json:"codemirrorMimeType"`
	CodemirrorMode     string   `json:"codemirrorMode"`
	Color              string   `json:"color"`
	Extensions         []string `json:"extensions"`
	Name               string   `json:"name"`
	TmScope            string   `json:"tmScope"`
	Wrap               bool     `json:"wrap"`
}

type PasteLanguageStats struct {
	Language   LanguageStats `json:"language"`
	Percentage float64       `json:"percentage"`
}

type CompactPasteHistory struct {
	EditedAt time.Time `json:"editedAt"`
	ID       string    `json:"id"`
}

func (c *Client) GetPaste(ctx context.Context, pasteID string) (*Paste, error) {
	url := fmt.Sprintf("%s/pastes/%s", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http Request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err == nil {
			return nil, fmt.Errorf("API Error (%s): %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("API returned non-200 status: %s, (body %s)",
				res.Status, string(bodyBytes))
		}
	}

	var paste Paste
	if err := json.NewDecoder(res.Body).Decode(&paste); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %w", err)
	}

	return &paste, nil
}

func (c *Client) GetPasteStats(ctx context.Context, pasteID string) (*Stats, error) {
	url := fmt.Sprintf("%s/pastes/%s/stats", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http request failed : %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err == nil {
			return nil, fmt.Errorf("API Error (%s): %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("API returned non-200 status: %s, (body %s)",
				res.Status, string(bodyBytes))
		}
	}

	var stats Stats
	if err := json.NewDecoder(res.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %w", err)
	}

	return &stats, nil
}

func (c *Client) CreatePaste(ctx context.Context, options CreatePasteOptions) (*Paste, error) {
	if len(options.Pasties) == 0 {
		return nil, fmt.Errorf("at least one pasty should be present")
	}

	// Converting the struct into JSON data
	jsonData, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("could not marshal paste options: %w", err)
	}

	// Creating the request
	url := fmt.Sprintf("%s/pastes", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
	}

	// Executing the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("api error (%s): %s", res.Status, apiError.StatusMessage)
		}
		return nil, fmt.Errorf("api returned non-201 status: %s", res.Status)
	}

	var newPaste Paste
	if err := json.NewDecoder(res.Body).Decode(&newPaste); err != nil {
		return nil, fmt.Errorf("could not decode json response: %w", err)
	}

	return &newPaste, nil
}

func (c *Client) GetPasteLanguageStats(ctx context.Context, pasteID string) ([]PasteLanguageStats, error) {
	url := fmt.Sprintf("%s/pastes/%s/langs", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err == nil {
			return nil, fmt.Errorf("API Error (%s) : %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("API returned non-200 status: %s, (body %s)",
				res.Status, string(bodyBytes))
		}
	}

	var pasteLangStats []PasteLanguageStats

	if err := json.NewDecoder(res.Body).Decode(&pasteLangStats); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %w", err)
	}

	return pasteLangStats, nil
}

func (c *Client) GetCompactPasteHistory(ctx context.Context, pasteID string) ([]CompactPasteHistory, error) {
	url := fmt.Sprintf("%s/pastes/%s/history_compact", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err == nil {
			return nil, fmt.Errorf("API Error (%s) : %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("API returned non-200 status: %s, (body %s)",
				res.Status, string(bodyBytes))
		}
	}

	var compactPasteHistory []CompactPasteHistory
	if err := json.NewDecoder(res.Body).Decode(&compactPasteHistory); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %w", err)
	}

	return compactPasteHistory, nil
}
