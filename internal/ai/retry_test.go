package ai

import (
	"errors"
	"testing"
	"time"
)

func TestCalculateBackoff(t *testing.T) {
	tests := []struct {
		name     string
		retry    int
		expected time.Duration
	}{
		{"first retry", 0, 1 * time.Second},
		{"second retry", 1, 2 * time.Second},
		{"third retry", 2, 4 * time.Second},
		{"fourth retry", 3, 8 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBackoff(tt.retry)
			if result != tt.expected {
				t.Errorf("CalculateBackoff(%d) = %v, want %v", tt.retry, result, tt.expected)
			}
		})
	}
}

func TestRetryWithBackoff_Success(t *testing.T) {
	callCount := 0
	fn := func() error {
		callCount++
		return nil
	}

	err := RetryWithBackoff(fn, 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestRetryWithBackoff_FailOnce(t *testing.T) {
	callCount := 0
	fn := func() error {
		callCount++
		if callCount == 1 {
			return errors.New("temporary error")
		}
		return nil
	}

	err := RetryWithBackoff(fn, 3)
	if err != nil {
		t.Errorf("Expected no error after retry, got %v", err)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestRetryWithBackoff_AlwaysFails(t *testing.T) {
	callCount := 0
	testError := errors.New("persistent error")
	fn := func() error {
		callCount++
		return testError
	}

	err := RetryWithBackoff(fn, 3)
	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Should call initial attempt + 3 retries = 4 total
	if callCount != 4 {
		t.Errorf("Expected 4 calls (1 initial + 3 retries), got %d", callCount)
	}

	if !errors.Is(err, testError) {
		t.Errorf("Expected wrapped testError, got %v", err)
	}
}

func TestRetryWithBackoff_MaxRetries(t *testing.T) {
	callCount := 0
	fn := func() error {
		callCount++
		return errors.New("error")
	}

	maxRetries := 2
	RetryWithBackoff(fn, maxRetries)

	// Should call initial attempt + maxRetries
	expectedCalls := maxRetries + 1
	if callCount != expectedCalls {
		t.Errorf("Expected %d calls, got %d", expectedCalls, callCount)
	}
}
