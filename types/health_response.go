package types

import (
	"net/http"
	"strings"
	"time"
)

type HealthStatus string

const (
	HealthStatusOK                 HealthStatus = "OK"
	HealthStatusDown               HealthStatus = "DOWN"
	HealthStatusNotFound           HealthStatus = "NOT_FOUND"
	HealthStatusPartiallyAvailable HealthStatus = "PARTIALLY_AVAILABLE"
)

func (hs HealthStatus) String() string {
	return string(hs)
}

type HealthResponse struct {
	Status       string        `json:"status"`
	StatusCode   int           `json:"statusCode"`
	Source       string        `json:"source"`
	ResponseTime time.Duration `json:"responseTime"`
	SentAt       time.Time     `json:"sentAt"`
	Components   []Component   `json:"components"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func (hr HealthResponse) HealthStatus() HealthStatus {
	switch {
	case hr.StatusCode == http.StatusOK && strings.EqualFold(hr.Status, HealthStatusOK.String()):
		return HealthStatusOK
	case hr.StatusCode == http.StatusOK && strings.EqualFold(hr.Status, HealthStatusPartiallyAvailable.String()):
		return HealthStatusPartiallyAvailable
	case componentsAreDown(hr.Components):
		return HealthStatusPartiallyAvailable
	case hr.StatusCode == http.StatusNotFound:
		return HealthStatusNotFound
	default:
		return HealthStatusDown
	}
}

func componentsAreDown(components []Component) bool {
	for _, component := range components {
		if !strings.EqualFold(component.Status, HealthStatusOK.String()) {
			return true
		}
	}
	return false
}
