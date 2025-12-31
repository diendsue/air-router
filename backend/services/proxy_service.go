package services

import (
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"net/http"
	"sync/atomic"

	"air_router/cache"
	"air_router/models"
	"air_router/utils"

	"github.com/gin-gonic/gin"
)

// Global counter for randomized account selection
var globalAccountCounter uint64

// Threshold for resetting the counter to avoid overflow
const counterResetThreshold = (1 << 63) - 100000

func init() {
	// Initialize with a secure random number between 10w and 20w
	n, _ := rand.Int(rand.Reader, big.NewInt(100000))
	globalAccountCounter = 100000 + n.Uint64()
}

// ProxyService handles proxy request routing and retry logic
type ProxyService struct {
	HTTPClient *http.Client
}

// NewProxyService creates a new ProxyService
func NewProxyService() *ProxyService {
	return &ProxyService{
		HTTPClient: utils.HTTPClient,
	}
}

// TryWithAccount attempts to forward request to a specific account
func (s *ProxyService) TryWithAccount(c *gin.Context, account models.Account, path string, bodyBytes []byte, headers http.Header) (*http.Response, bool, []byte) {
	targetURL := utils.BuildTargetURL(account, path)

	req, err := utils.CreateProxyRequest(c.Request.Method, targetURL, bodyBytes, account, headers)
	if err != nil {
		return nil, false, nil
	}

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, false, nil
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return resp, false, bodyBytes
	}

	return resp, true, nil
}

// TryWithRetryModel attempts to forward request using accounts that support the model
// Returns (success, lastResponse, lastResponseBody)
func (s *ProxyService) TryWithRetryModel(c *gin.Context, path string, modelID string, bodyBytes []byte) (bool, *http.Response, []byte) {
	accounts := cache.GetAccountsForModel(modelID)
	if len(accounts) == 0 {
		return false, nil, nil
	}

	log.Printf("[ProxyService] Model: %s, Accounts: %d", modelID, len(accounts))

	// Retry at most 2 times
	maxAttempts := 2
	if len(accounts) < 2 {
		maxAttempts = len(accounts)
	}

	var lastResp *http.Response
	var lastRespBody []byte

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Atomically increment and get current value
		current := atomic.AddUint64(&globalAccountCounter, 1)

		// Reset if approaching int64 limit to avoid overflow
		if current >= counterResetThreshold {
			n, _ := rand.Int(rand.Reader, big.NewInt(100000))
			atomic.StoreUint64(&globalAccountCounter, 100000+n.Uint64())
			current = globalAccountCounter
		}

		accountIndex := current % uint64(len(accounts))
		account := accounts[accountIndex]

		log.Printf("[ProxyService] Attempt %d/%d with account %s (ID: %d)", attempt+1, maxAttempts, account.Name, account.ID)

		resp, success, respBody := s.TryWithAccount(c, account, path, bodyBytes, c.Request.Header)
		if resp != nil {
			// Keep track of last response
			lastResp = resp
			lastRespBody = respBody

			defer resp.Body.Close()
			if success {
				// Stream response
				utils.StreamResponse(c, resp)
				log.Printf("[ProxyService] Success with account %s (ID: %d)", account.Name, account.ID)
				return true, nil, nil
			}
		}
	}

	return false, lastResp, lastRespBody
}
