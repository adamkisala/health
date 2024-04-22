package healthchecker_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/adamkisala/weaviate-health/internal/healthchecker"
)

func TestRunner_Run(t *testing.T) {
	type fields struct {
		logger            *slog.Logger
		requester         healthchecker.Requester
		responseProcessor healthchecker.ResponseProcessor
		sourcesProvider   healthchecker.SourcesProvider
		workers           int64
		checkInterval     time.Duration
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := healthchecker.NewRunner(healthchecker.RunnerParams{
				Logger:            tt.fields.logger,
				Requester:         tt.fields.requester,
				ResponseProcessor: tt.fields.responseProcessor,
				SourcesProvider:   tt.fields.sourcesProvider,
				Workers:           tt.fields.workers,
				CheckInterval:     tt.fields.checkInterval,
			})
			if err := r.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
