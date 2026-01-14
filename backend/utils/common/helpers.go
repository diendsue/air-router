package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"air_router/constants"
	"air_router/models"

	"github.com/gin-gonic/gin"
)

// SendJSONResponse sends a JSON response with the specified status code
func SendJSONResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(c *gin.Context, statusCode int, message, errType string) {
	SendAPIError(c, statusCode, message, errType)
}

// ParseIDParam parses and validates an ID parameter from the URL
func ParseIDParam(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID parameter")
	}
	return id, nil
}

// ValidateModelProvider validates if the provider is supported
func ValidateModelProvider(provider models.Provider) bool {
	validProviders := map[models.Provider]bool{
		models.ProviderChat:   true,
		models.ProviderClaude: true,
		models.ProviderCodex:  true,
		models.ProviderGemini: true,
	}
	return validProviders[provider]
}

// GetEnvOrDefault gets an environment variable or returns a default value
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ExtractModelID extracts model ID from request body
func ExtractModelID(bodyBytes []byte) string {
	var requestBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
		return ""
	}
	if modelID, ok := requestBody["model"].(string); ok {
		return modelID
	}
	return ""
}

// Global counter for randomized selection
var globalCounter uint64

func init() {
	// Initialize with a secure random number between 100000 and 200000
	n, _ := rand.Int(rand.Reader, big.NewInt(100000))
	globalCounter = 100000 + n.Uint64()
}

// GetGlobalCounter returns the global counter for external access
func GetGlobalCounter() *uint64 {
	return &globalCounter
}

// GetRandomIndex selects a random index using global counter
func GetRandomIndex(length int) int {
	if length <= 0 {
		return 0
	}

	if length == 1 {
		return 0
	}

	// Atomically increment and get current value
	counter := GetGlobalCounter()
	current := atomic.AddUint64(counter, 1)

	// Reset if approaching int64 limit to avoid overflow
	if current >= constants.CounterResetThreshold {
		n, _ := rand.Int(rand.Reader, big.NewInt(100000))
		atomic.StoreUint64(counter, 100000+n.Uint64())
		current = *counter
	}

	index := int(current % uint64(length))
	return index
}

// GetRandomElement selects a random element from a slice
func GetRandomElement[T any](elements []T) T {
	if len(elements) == 0 {
		var zero T
		return zero
	}

	if len(elements) == 1 {
		return elements[0]
	}

	index := GetRandomIndex(len(elements))
	return elements[index]
}

// GetCurrentTimestamp returns the current timestamp in milliseconds
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// ResetCounterIfNeeded checks and resets the global counter if needed
func ResetCounterIfNeeded() {
	counter := GetGlobalCounter()
	current := atomic.LoadUint64(counter)

	// Reset if approaching int64 limit to avoid overflow
	if current >= constants.CounterResetThreshold {
		n, _ := rand.Int(rand.Reader, big.NewInt(100000))
		atomic.StoreUint64(counter, 100000+n.Uint64())
		log.Printf("[Counter] Reset global counter to: %d", 100000+n.Uint64())
	}
}
