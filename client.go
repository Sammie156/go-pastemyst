package gopastemyst

import (
	"net/http"
	"time"
)

// Current endpoint for V3 API
const BaseURL = "https://beta.myst.rs/api/v3"

type Client struct {
	baseURL    string
	apiToken   string // Storing the user's API Token
	httpClient *http.Client
}

func NewClient(apiToken string) *Client {
	return &Client{
		baseURL:  BaseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type APIError struct {
	StatusMessage string `json:"statusMessage"`
}

func (e APIError) Error() string {
	return e.StatusMessage
}
