package virustotal

import (
	"io"
	"net/http"
)

// Client holds the configuration for the API client
type Client struct {
	apiKey string
	client *http.Client
}

// NewClient initializes and returns a new VirusTotal API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// makeRequest helps in making requests to the VirusTotal API
func (c *Client) makeRequest(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-apikey", c.apiKey)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.client.Do(req)
}
