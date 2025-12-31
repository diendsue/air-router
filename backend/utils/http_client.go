package utils

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"
)

// SkipTLSVerify determines whether to skip TLS certificate verification
var SkipTLSVerify = os.Getenv("SKIP_TLS_VERIFY") != "false" // Default: true

// HTTPClient is the shared HTTP client for the application
var HTTPClient *http.Client

func init() {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: SkipTLSVerify,
		},
	}
	HTTPClient = &http.Client{
		Transport: transport,
		Timeout:   0, // No timeout to support long connections and streaming
	}
}

// InitHTTPClient initializes the HTTP client
func InitHTTPClient() *http.Client {
	return HTTPClient
}

// GetHTTPClient returns the shared HTTP client
func GetHTTPClient() *http.Client {
	return HTTPClient
}
