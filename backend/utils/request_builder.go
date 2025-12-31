package utils

import (
	"bytes"
	"net/http"

	"air_router/models"
)

const defaultAnthropicVersion = "2023-06-01"
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

// CreateProxyRequest creates an HTTP request for proxying
// isClaude indicates whether this is a Claude API request
func CreateProxyRequest(method, targetURL string, bodyBytes []byte, account models.Account, headers http.Header, isClaude bool) (*http.Request, error) {
	req, err := http.NewRequest(method, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "identity")
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set API key based on API type
	if isClaude {
		req.Header.Set("X-Api-Key", account.APIKey)
		// Remove Authorization header if present
		req.Header.Del("Authorization")
		// Set or get anthropic-version header
		anthropicVersion := headers.Get("anthropic-version")
		if anthropicVersion == "" {
			anthropicVersion = defaultAnthropicVersion
		}
		req.Header.Set("anthropic-version", anthropicVersion)
	} else {
		req.Header.Set("Authorization", "Bearer "+account.APIKey)
	}

	// Always set User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", userAgent)
	}

	return req, nil
}
