package constants

const (
	// API Constants
	DefaultAnthropicVersion = "2023-06-01"
	DefaultUserAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

	// Pagination Constants
	DefaultPageSize = 10
	MaxPageSize     = 100

	// HTTP Constants
	StreamBufferSize = 4096

	// Cache Constants
	CounterResetThreshold = (1 << 63) - 100000
)
