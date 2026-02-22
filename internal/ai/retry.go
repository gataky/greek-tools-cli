package ai

import (
	"fmt"
	"time"
)

// CalculateBackoff returns the wait duration for exponential backoff
// Progression: 1s, 2s, 4s, 8s, etc.
func CalculateBackoff(retryCount int) time.Duration {
	if retryCount <= 0 {
		return time.Second
	}
	// 2^retryCount seconds
	return time.Duration(1<<uint(retryCount)) * time.Second
}

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// RetryWithBackoff executes a function with exponential backoff retry logic
// It will retry up to maxRetries times (total of maxRetries+1 attempts)
// Returns the final error if all attempts fail
func RetryWithBackoff(fn RetryableFunc, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Don't sleep after the last attempt
		if attempt < maxRetries {
			backoff := CalculateBackoff(attempt)
			time.Sleep(backoff)
		}
	}

	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
