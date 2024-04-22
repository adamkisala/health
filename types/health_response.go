package types

import (
	"time"
)

type HealthResponse struct {
	Status       string        `json:"status"`
	StatusCode   int           `json:"statusCode"`
	Source       string        `json:"source"`
	ResponseTime time.Duration `json:"responseTime"`
	TimeStamp    time.Time     `json:"timeStamp"`
	Components   []Component   `json:"components"`
	Error        error         `json:"error"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
