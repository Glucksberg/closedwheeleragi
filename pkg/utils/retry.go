// Package utils provides common utility functions for the AGI agent.
package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// RetryConfig defines the configuration for retry logic
type RetryConfig struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	JitterFactor float64
}

// DefaultRetryConfig returns a standard retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     10 * time.Second,
		JitterFactor: 0.2,
	}
}

// ExecuteWithRetry executes a function with exponential backoff and jitter
func ExecuteWithRetry(operation func() error, config RetryConfig) error {
	var lastErr error
	delay := config.InitialDelay

	for i := 0; i <= config.MaxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err
		if i == config.MaxRetries {
			break
		}

		// Calculate jittered delay
		jitter := float64(delay) * config.JitterFactor
		actualDelay := delay + time.Duration((rand.Float64()*2-1)*jitter)

		time.Sleep(actualDelay)

		// Exponential backoff
		delay *= 2
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", config.MaxRetries, lastErr)
}

// IsRetryableError determines if an HTTP status code is retryable (429 or 5xx)
func IsRetryableError(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || (statusCode >= 500 && statusCode <= 599)
}
