package healthchecker

import (
	"context"
	"net/url"

	"github.com/adamkisala/health/types"
)

type MockRequester struct {
	RequestErr      error
	RequestRes      map[string]types.HealthResponse
	requestsSources []*url.URL
}

func (m *MockRequester) Request(ctx context.Context, source *url.URL) (types.HealthResponse, error) {
	if m.RequestRes == nil {
		m.RequestRes = make(map[string]types.HealthResponse)
	}
	m.requestsSources = append(m.requestsSources, source)
	return m.RequestRes[source.String()], m.RequestErr
}

func (m *MockRequester) RequestsSources() []*url.URL {
	return m.requestsSources
}
