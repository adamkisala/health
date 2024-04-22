package healthchecker

import (
	"context"

	"github.com/adamkisala/weaviate-health/types"
)

type MockResponseProcessor struct {
	ProcessErr         error
	processedResponses []types.HealthResponse
}

func (m *MockResponseProcessor) Process(ctx context.Context, response types.HealthResponse) error {
	m.processedResponses = append(m.processedResponses, response)
	return m.ProcessErr
}

func (m *MockResponseProcessor) ProcessedResponses() []types.HealthResponse {
	return m.processedResponses
}
