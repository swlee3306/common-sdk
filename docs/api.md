# Common-SDK API Documentation

## Overview
Common-SDK is a Go library for multicast communication with compression, encryption, and monitoring capabilities.

## Features
- Message compression (Gzip, LZ4)
- Message encryption (AES-256)
- Prometheus metrics
- Structured logging
- Health checks
- Retry logic with backoff strategies
- Connection pooling

## Installation
```bash
go get github.com/your-org/common-sdk
```

## Quick Start

### Basic Usage
```go
package main

import (
    "fmt"
    "github.com/your-org/common-sdk/compression"
    "github.com/your-org/common-sdk/encryption"
)

func main() {
    // Compression
    compressor := compression.NewCompressor(compression.Gzip)
    compressed, err := compressor.Compress([]byte("Hello World"))
    if err != nil {
        panic(err)
    }
    
    // Encryption
    encryptor, err := encryption.NewEncryptor("your-secret-key")
    if err != nil {
        panic(err)
    }
    
    encrypted, err := encryptor.EncryptBase64("Hello World")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Compressed: %v\n", compressed)
    fmt.Printf("Encrypted: %s\n", encrypted)
}
```

## API Reference

### Compression
- `NewCompressor(algorithm Algorithm) *Compressor`
- `Compress(data []byte) ([]byte, error)`
- `Decompress(data []byte) ([]byte, error)`

### Encryption
- `NewEncryptor(key string) (*Encryptor, error)`
- `Encrypt(plaintext []byte) ([]byte, error)`
- `Decrypt(ciphertext []byte) ([]byte, error)`
- `EncryptBase64(plaintext string) (string, error)`
- `DecryptBase64(encrypted string) (string, error)`

### Metrics
- `NewMetrics() *Metrics`
- `RecordMessageSent(size int)`
- `RecordMessageReceived(size int)`
- `RecordProcessingTime(duration time.Duration)`
- `SetActiveConnections(count int)`
- `RecordError()`

### Logging
- `NewLogger(level LogLevel) *Logger`
- `WithField(key string, value interface{}) *Logger`
- `Debug(message string, fields ...map[string]interface{})`
- `Info(message string, fields ...map[string]interface{})`
- `Warn(message string, fields ...map[string]interface{})`
- `Error(message string, fields ...map[string]interface{})`

### Health Checks
- `NewHealthChecker(version string) *HealthChecker`
- `AddCheck(name string, checkFunc func() Check)`
- `GetHealth() HealthCheck`

### Retry Logic
- `NewRetryer(config RetryConfig) *Retryer`
- `Execute(ctx context.Context, operation func() error) error`
- `ExecuteWithResult(ctx context.Context, operation func() (interface{}, error)) (interface{}, error)`

## Examples

### Complete Example
```go
package main

import (
    "context"
    "time"
    "github.com/your-org/common-sdk/compression"
    "github.com/your-org/common-sdk/encryption"
    "github.com/your-org/common-sdk/metrics"
    "github.com/your-org/common-sdk/logging"
    "github.com/your-org/common-sdk/health"
    "github.com/your-org/common-sdk/retry"
)

func main() {
    // Initialize components
    logger := logging.NewLogger(logging.INFO)
    metrics := metrics.NewMetrics()
    healthChecker := health.NewHealthChecker("1.0.0")
    retryer := retry.NewRetryer(retry.RetryConfig{
        MaxAttempts: 3,
        BaseDelay:   time.Second,
        MaxDelay:    time.Minute,
        Strategy:    retry.Exponential,
        Multiplier:  2.0,
    })
    
    // Add health checks
    healthChecker.AddCheck("compression", func() health.Check {
        return health.Check{Status: health.Healthy, Message: "OK"}
    })
    
    // Use retry logic
    ctx := context.Background()
    err := retryer.Execute(ctx, func() error {
        // Your operation here
        return nil
    })
    
    if err != nil {
        logger.Error("Operation failed", map[string]interface{}{
            "error": err.Error(),
        })
    }
}
```

## Configuration

### Environment Variables
- `LOG_LEVEL`: Logging level (DEBUG, INFO, WARN, ERROR)
- `METRICS_ENABLED`: Enable metrics collection
- `HEALTH_CHECK_INTERVAL`: Health check interval

## License
MIT License
