package errors

import (
	"context"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay time.Duration
	Multiplier float64
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 3,
	InitialDelay: 1 * time.Second,
	MaxDelay: 30 * time.Second,
	Multiplier: 2.0,
}

func RetryWithBackoff(ctx context.Context, config RetryConfig, fn func() error) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return Wrap(ctx.Err(), ErrCodeTimeout, "retry cancelled")
			case <-time.After(delay):
			}

			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if !isRetryable(lastErr) {
			return lastErr
		}
	}

	return Wrap(lastErr, ErrCodeNetwork, "max retry attempts exceeded").WithDetail("attempts", config.MaxAttempts)
}

func isRetryable(err error) bool {
	code := GetCode(err)
	return code == ErrCodeNetwork || code == ErrCodeTimeout
}

func RetryOnError(ctx context.Context, fn func() error) error {
	return RetryWithBackoff(ctx, DefaultRetryConfig, fn)
}
