package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/adamkisala/weaviate-health/types"
	"github.com/pkg/errors"
	httpstat "github.com/tcnksm/go-httpstat"
)

type HealthResponse struct {
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

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPRequester struct {
	doer                   Doer
	defaultTimeout         time.Duration
	defaultHealthCheckPath string
}

type HTTPRequesterParams struct {
	Doer                   Doer
	DefaultTimeout         time.Duration
	DefaultHealthCheckPath string
}

func NewHTTPRequester(params HTTPRequesterParams) *HTTPRequester {
	return &HTTPRequester{
		doer:                   params.Doer,
		defaultTimeout:         params.DefaultTimeout,
		defaultHealthCheckPath: params.DefaultHealthCheckPath,
	}
}

func (hr *HTTPRequester) Request(ctx context.Context, source *url.URL) (types.HealthResponse, error) {
	healthSourceURL, err := url.JoinPath(source.Host, hr.defaultHealthCheckPath)
	if err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to join URL path")
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, hr.defaultTimeout)
	defer cancel()

	var statResults httpstat.Result
	req, err := http.NewRequestWithContext(
		httpstat.WithHTTPStat(ctxWithTimeout, &statResults),
		http.MethodGet,
		healthSourceURL,
		nil)
	if err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to create request")
	}
	resp, err := hr.doer.Do(req)
	if errors.Is(err, context.DeadlineExceeded) {
		return types.HealthResponse{
			StatusCode:   http.StatusRequestTimeout,
			Status:       http.StatusText(http.StatusRequestTimeout),
			ResponseTime: statResults.ServerProcessing,
			Source:       source.Host,
		}, nil
	}
	if err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()
	var healthResponse HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to decode response")
	}
	return types.HealthResponse{
		StatusCode:   resp.StatusCode,
		Status:       resp.Status,
		Source:       source.Host,
		ResponseTime: statResults.ServerProcessing,
	}, nil
}