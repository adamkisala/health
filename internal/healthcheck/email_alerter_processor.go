package healthcheck

import (
	"context"

	"github.com/adamkisala/weaviate-health/types"
)

type EmailAlerterProcessor struct {
}

type EmailAlerterProcessorParams struct {
}

func NewEmailAlerterProcessor(params EmailAlerterProcessorParams) *EmailAlerterProcessor {
	return &EmailAlerterProcessor{}
}

func (p EmailAlerterProcessor) Process(ctx context.Context, response types.HealthResponse) error {
	// TODO: implement email alerting logic if response.Failed() for example
	// consolidate knowledge about the health of response to one component
	return nil
}
