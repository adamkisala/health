package healthchecker_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/adamkisala/weaviate-health/internal/healthchecker"
	"github.com/adamkisala/weaviate-health/tests"
	mocks "github.com/adamkisala/weaviate-health/tests/mocks/healthchecker"
	"github.com/adamkisala/weaviate-health/types"
	"github.com/stretchr/testify/require"
)

func TestRunner_Run(t *testing.T) {
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()
	ctxWithTimeout, ctxWithTimeoutCncl := context.WithTimeout(context.Background(), 2*time.Second)
	defer ctxWithTimeoutCncl()
	type fields struct {
		requester         healthchecker.Requester
		responseProcessor *mocks.MockResponseProcessor
		sourcesProvider   healthchecker.SourcesProvider
		workers           int64
		checkInterval     time.Duration
	}
	type args struct {
		ctx context.Context
	}
	testCases := []struct {
		name                           string
		fields                         fields
		args                           args
		timeToProcess                  time.Duration
		expectedLoggerOutput           string
		expectedProcessedResponsesSize int
	}{
		{
			name: "runner is cancellable",
			args: args{ctx: cancelledCtx},
			fields: fields{
				checkInterval:     1 * time.Second,
				responseProcessor: &mocks.MockResponseProcessor{},
			},
			timeToProcess:        200 * time.Millisecond,
			expectedLoggerOutput: "context done - stopping healthcheck runner",
		},
		{
			name: "runner is able to process messages without splitting work",
			args: args{ctx: ctxWithTimeout},
			fields: fields{
				checkInterval: 500 * time.Millisecond,
				sourcesProvider: &mocks.MockSourcesProvider{
					ProvideResult: []*url.URL{
						tests.MustParseURL(t, "https://valid-url-a.com/health"),
						tests.MustParseURL(t, "https://valid-url-b.com/health"),
					},
				},
				requester: &mocks.MockRequester{
					RequestRes: map[string]types.HealthResponse{
						"https://valid-url-a.com/health": {
							Status:     "OK",
							StatusCode: http.StatusOK,
						},
						"https://valid-url-b.com/health": {
							Status:     "Not Found",
							StatusCode: http.StatusNotFound,
						},
					},
				},
				responseProcessor: &mocks.MockResponseProcessor{},
			},
			timeToProcess:                  10 * time.Second,
			expectedLoggerOutput:           "completed healthcheck run",
			expectedProcessedResponsesSize: 2,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			loggerHandler := slog.NewTextHandler(&buf, nil)
			loggerWithHandler := slog.New(loggerHandler)
			r := healthchecker.NewRunner(healthchecker.RunnerParams{
				Logger:            loggerWithHandler,
				Requester:         tt.fields.requester,
				ResponseProcessor: tt.fields.responseProcessor,
				SourcesProvider:   tt.fields.sourcesProvider,
				Workers:           tt.fields.workers,
				CheckInterval:     tt.fields.checkInterval,
			})

			go func() {
				select {
				case <-tt.args.ctx.Done():
					return
				case <-time.After(tt.timeToProcess):
					t.Errorf("Runner.Run() did not return after %v", tt.timeToProcess)
					return
				}
			}()

			err := r.Run(tt.args.ctx)
			require.NoError(t, err)
			require.Contains(t, buf.String(), tt.expectedLoggerOutput)
			require.GreaterOrEqual(t, len(tt.fields.responseProcessor.ProcessedResponses()), tt.expectedProcessedResponsesSize)
		})
	}
}
