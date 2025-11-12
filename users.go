package gopastemyst

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Check : https://docs.beta.myst.rs/users

func (c *Client) GetUser(ctx context.Context, username string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", c.baseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not request http error: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("API error (%s): %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("API returned non-200 status: %s, body(%s)", res.Status, string(bodyBytes))
		}
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("could not decode JSON response: %w", err)
	}

	return &user, nil
}

// func (c *Client)
