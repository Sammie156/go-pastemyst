package gopastemyst

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Currently only checking some fields.

// TODO: Add more fields based on the responses
// Check https://docs.beta.myst.rs/pastes

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

// TODO: Continue from here tomorrow and make a working
//
//	Struct for each Paste
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
