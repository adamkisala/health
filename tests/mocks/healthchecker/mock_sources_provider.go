package healthchecker

import (
	"context"
	"net/url"
)

type MockSourcesProvider struct {
	ProvideResult []*url.URL
	ProvideError  error
}

func (m *MockSourcesProvider) Provide(ctx context.Context) ([]*url.URL, error) {
	return m.ProvideResult, m.ProvideError
}
