package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/adamkisala/health/types"
	"github.com/pkg/errors"
	httpstat "github.com/tcnksm/go-httpstat"
)

type HealthResponse struct {
	Status     string      `json:"status"`
	Components []Component `json:"components"`
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
	defaultHealthCheckPath string
}

type HTTPRequesterParams struct {
	Doer                   Doer
	DefaultHealthCheckPath string
}

func NewHTTPRequester(params HTTPRequesterParams) *HTTPRequester {
	return &HTTPRequester{
		doer:                   params.Doer,
		defaultHealthCheckPath: params.DefaultHealthCheckPath,
	}
}

func (hr *HTTPRequester) Request(ctx context.Context, source *url.URL) (types.HealthResponse, error) {
	healthSourceURL, err := url.JoinPath(source.String(), hr.defaultHealthCheckPath)
	if err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to join URL path")
	}
	var statResults httpstat.Result
	req, err := http.NewRequestWithContext(
		httpstat.WithHTTPStat(ctx, &statResults),
		http.MethodGet,
		healthSourceURL,
		nil)
	if err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to create request")
	}
	sentAt := time.Now().UTC()
	resp, err := hr.doer.Do(req)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTransientStatusCode) {
		if resp == nil {
			return types.HealthResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     http.StatusText(http.StatusInternalServerError),
				Source:     source.String(),
				SentAt:     sentAt,
			}, nil
		}
		return types.HealthResponse{
			StatusCode:   resp.StatusCode,
			Status:       http.StatusText(resp.StatusCode),
			ResponseTime: statResults.ServerProcessing,
			Source:       source.String(),
			SentAt:       sentAt,
		}, nil
	}
	if err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to make request")
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	var healthResponse HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		return types.HealthResponse{}, errors.Wrap(err, "failed to decode response")
	}

	return types.HealthResponse{
		Status:       healthResponse.Status,
		StatusCode:   resp.StatusCode,
		Source:       source.String(),
		ResponseTime: statResults.ServerProcessing,
		SentAt:       sentAt,
		Components:   responseComponentsToDomain(healthResponse.Components),
	}, nil
}

func responseComponentsToDomain(components []Component) []types.Component {
	var domainComponents []types.Component
	for _, component := range components {
		domainComponents = append(domainComponents, types.Component{
			Name:   component.Name,
			Status: component.Status,
		})
	}
	return domainComponents
}
