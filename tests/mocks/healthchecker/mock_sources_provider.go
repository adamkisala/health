package healthchecker

import (
	"context"

	"github.com/adamkisala/weaviate-health/types"
)

type MockSourcesProvider struct {
	ProvideResult types.Sources
	ProvideError  error
}

func (m *MockSourcesProvider) Provide(ctx context.Context) (types.Sources, error) {
	return m.ProvideResult, m.ProvideError
}
