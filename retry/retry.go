package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

type BackoffStrategy string

const (
	Linear     BackoffStrategy = "linear"
	Exponential BackoffStrategy = "exponential"
	Fixed      BackoffStrategy = "fixed"
)

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Strategy    BackoffStrategy
	Multiplier  float64
}

type Retryer struct {
	config RetryConfig
}

func NewRetryer(config RetryConfig) *Retryer {
	return &Retryer{config: config}
}

func (r *Retryer) Execute(ctx context.Context, operation func() error) error {
	var lastErr error
	
	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		err := operation()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		if attempt == r.config.MaxAttempts {
			break
		}
		
		delay := r.calculateDelay(attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	
	return fmt.Errorf("operation failed after %d attempts: %w", r.config.MaxAttempts, lastErr)
}

func (r *Retryer) calculateDelay(attempt int) time.Duration {
	var delay time.Duration
	
	switch r.config.Strategy {
	case Linear:
		delay = r.config.BaseDelay * time.Duration(attempt)
	case Exponential:
		delay = time.Duration(float64(r.config.BaseDelay) * math.Pow(r.config.Multiplier, float64(attempt-1)))
	case Fixed:
		delay = r.config.BaseDelay
	default:
		delay = r.config.BaseDelay
	}
	
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}
	
	return delay
}

func (r *Retryer) ExecuteWithResult(ctx context.Context, operation func() (interface{}, error)) (interface{}, error) {
	var lastErr error
	
	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		result, err := operation()
		if err == nil {
			return result, nil
		}
		
		lastErr = err
		
		if attempt == r.config.MaxAttempts {
			break
		}
		
		delay := r.calculateDelay(attempt)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	
	return nil, fmt.Errorf("operation failed after %d attempts: %w", r.config.MaxAttempts, lastErr)
}
