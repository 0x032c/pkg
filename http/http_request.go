package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPRequest executes an HTTP request with given parameters and decodes the response as JSON into responseStruct.
// method: "GET", "POST", etc.
// headers: key-value map of request headers.
// queryParams: key-value map of URL query parameters.
// body: request body (will be marshaled to JSON if not nil).
// responseStruct: pointer to struct to decode JSON response into.
// timeout: timeout for HTTP request (default 10s if <=0).
// Returns error if the request or decoding fails.
func HTTPRequest(
	ctx context.Context,
	method string,
	requestURL string,
	headers map[string]string,
	queryParams map[string]string,
	body interface{},
	responseStruct interface{},
	timeout time.Duration,
) error {
	// Parse URL and add query parameters
	urlObj, err := url.Parse(requestURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}
	query := urlObj.Query()
	for key, value := range queryParams {
		query.Add(key, value)
	}
	urlObj.RawQuery = query.Encode()

	// Prepare request body
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		reqBody = bytes.NewBuffer(data)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, method, urlObj.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Content-Type if not provided
	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers["Content-Type"]; !ok && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set up HTTP client with timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	client := &http.Client{Timeout: timeout}

	// Do request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Accept 2xx as success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received non-2xx response: %s, body: %s", resp.Status, string(bodyBytes))
	}

	// Decode JSON response if responseStruct is not nil
	if responseStruct != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, responseStruct); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
