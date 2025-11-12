package gopastemyst

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Check https://docs.beta.myst.rs/pastes

// TODO: Add documentation to each struct type and function

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

func (c *Client) GetPasteAtSpecificEdit(ctx context.Context, pasteID string, historyID string) (*Paste, error) {
	url := fmt.Sprintf("%s/pastes/%s/history/%s", c.baseURL, pasteID, historyID)

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

	var paste Paste

	if err := json.NewDecoder(res.Body).Decode(&paste); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return &paste, nil
}

func (c *Client) GetDiffAtCertainEdit(ctx context.Context, pasteID string, historyID string) (*PasteDiff, error) {
	url := fmt.Sprintf("%s/pastes/%s/history/%s/diff", c.baseURL, pasteID, historyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %w", err)
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

	var pasteDiff PasteDiff
	if err := json.NewDecoder(res.Body).Decode(&pasteDiff); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return &pasteDiff, nil
}

func (c *Client) DownloadPasteAsZip(ctx context.Context, pasteID string) ([]byte, error) {
	url := fmt.Sprintf("%s/pastes/%s.zip", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %w", err)
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

	zipData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return zipData, nil
}

func (c *Client) IsPasteEncrypted(ctx context.Context, pasteID string) (bool, error) {
	url := fmt.Sprintf("%s/pastes/%s/encrypted", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("could not create http request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err == nil {
			return false, fmt.Errorf("API Error (%s) : %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return false, fmt.Errorf("API returned non-200 status: %s, (body %s)",
				res.Status, string(bodyBytes))
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	isEncrypted, err := strconv.ParseBool(string(bodyBytes))
	if err != nil {
		return false, fmt.Errorf("api returned non-boolean value: %s", string(bodyBytes))
	}

	return isEncrypted, nil
}

func (c *Client) IsPasteStarred(ctx context.Context, pasteID string) (bool, error) {
	url := fmt.Sprintf("%s/pastes/%s/star", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("could not create http request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err == nil {
			return false, fmt.Errorf("API Error (%s) : %s", res.Status, apiError.StatusMessage)
		}

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return false, fmt.Errorf("API returned non-200 status: %s, (body %s)",
				res.Status, string(bodyBytes))
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	isStarred, err := strconv.ParseBool(string(bodyBytes))
	fmt.Printf("Body Bytes: %s\n", string(bodyBytes))
	if err != nil {
		return false, fmt.Errorf("api returned non-boolean value: %s", string(bodyBytes))
	}

	return isStarred, nil
}

// TODO: Get back to this and finish all Paste endpoints

func (c *Client) StarPaste(ctx context.Context, pasteID string) error {
	url := fmt.Sprintf("%s/pastes/%s/star", c.baseURL, pasteID)

	if c.apiToken == "" {
		return fmt.Errorf("please provide API token")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError
		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("api error (%s): %s", res.Status, apiError.StatusMessage)
		}

		return fmt.Errorf("api returned non-204 status : %s", res.Status)
	}

	return nil
}

func (c *Client) PinPaste(ctx context.Context, pasteID string) error {
	url := fmt.Sprintf("%s/pastes/%s/pin", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		return fmt.Errorf("could not create http request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("api error (%s) : %s", res.Status, apiError.StatusMessage)
		}

		return fmt.Errorf("api error (%s): %s", res.Status, apiError.StatusMessage)
	}

	return nil
}

func (c *Client) PrivatePaste(ctx context.Context, pasteID string) error {
	url := fmt.Sprintf("%s/pastes/%s/private", c.baseURL, pasteID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(""))
	if err != nil {
		return fmt.Errorf("could not create http request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("api error (%s) : %s", res.Status, apiError.StatusMessage)
		}

		return fmt.Errorf("api error (%s): %s", res.Status, apiError.StatusMessage)
	}

	return nil
}

func (c *Client) EditPaste(ctx context.Context, pasteID string, options EditPasteOptions) (*Paste, error) {
	url := fmt.Sprintf("%s/pastes/%s", c.baseURL, pasteID)

	if c.apiToken == "" {
		return nil, fmt.Errorf("please provide API token")
	}

	jsonData, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("could not marshal edit options: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiError APIError

		if err := json.NewDecoder(res.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("API error (%s): %s", res.Status, apiError.StatusMessage)
		}

		return nil, fmt.Errorf("API error (%s): %s", res.Status, apiError.StatusMessage)
	}

	var paste Paste
	if err := json.NewDecoder(res.Body).Decode(&paste); err != nil {
		return nil, fmt.Errorf("could not decode json: %w", err)
	}

	return &paste, nil
}
