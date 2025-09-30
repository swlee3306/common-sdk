package health

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthStatus string

const (
	Healthy   HealthStatus = "healthy"
	Unhealthy HealthStatus = "unhealthy"
	Degraded  HealthStatus = "degraded"
)

type HealthCheck struct {
	Status    HealthStatus         `json:"status"`
	Timestamp time.Time           `json:"timestamp"`
	Version   string              `json:"version"`
	Checks    map[string]Check    `json:"checks"`
}

type Check struct {
	Status    HealthStatus `json:"status"`
	Message   string       `json:"message,omitempty"`
	Duration  string       `json:"duration,omitempty"`
}

type HealthChecker struct {
	version string
	checks  map[string]func() Check
}

func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		version: version,
		checks:  make(map[string]func() Check),
	}
}

func (hc *HealthChecker) AddCheck(name string, checkFunc func() Check) {
	hc.checks[name] = checkFunc
}

func (hc *HealthChecker) GetHealth() HealthCheck {
	start := time.Now()
	checks := make(map[string]Check)
	overallStatus := Healthy
	
	for name, checkFunc := range hc.checks {
		checkStart := time.Now()
		check := checkFunc()
		check.Duration = time.Since(checkStart).String()
		checks[name] = check
		
		if check.Status == Unhealthy {
			overallStatus = Unhealthy
		} else if check.Status == Degraded && overallStatus != Unhealthy {
			overallStatus = Degraded
		}
	}
	
	return HealthCheck{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Version:   hc.version,
		Checks:    checks,
	}
}

func (hc *HealthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	health := hc.GetHealth()
	
	w.Header().Set("Content-Type", "application/json")
	
	statusCode := http.StatusOK
	if health.Status == Unhealthy {
		statusCode = http.StatusServiceUnavailable
	} else if health.Status == Degraded {
		statusCode = http.StatusOK
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(health)
}
