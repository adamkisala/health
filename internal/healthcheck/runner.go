package healthcheck

import (
	"context"
	"log/slog"
	"net/url"
	"time"

	"github.com/adamkisala/weaviate-health/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// Requester is an interface for making health requests
type Requester interface {
	Request(ctx context.Context, source *url.URL) (types.HealthResponse, error)
}

type ResponseProcessor interface {
	Process(ctx context.Context, response types.HealthResponse) error
}

// SourcesProvider is an interface for providing sources
type SourcesProvider interface {
	Provide(ctx context.Context) (types.Sources, error)
}

type Runner struct {
	logger             *slog.Logger
	requester          Requester
	responseProcessors []ResponseProcessor
	sourcesProvider    SourcesProvider
	workers            int64
	checkInterval      time.Duration
}

type RunnerParams struct {
	Logger             *slog.Logger
	Requester          Requester
	ResponseProcessors []ResponseProcessor
	SourcesProvider    SourcesProvider
	Workers            int64
	CheckInterval      time.Duration
}

func NewRunner(params RunnerParams) *Runner {
	if params.Logger == nil {
		params.Logger = slog.Default()
	}
	if params.Workers == 0 {
		params.Workers = 1
	}
	return &Runner{
		logger: params.Logger.With(
			slog.String("package", "healthchecker"),
			slog.String("struct", "Runner"),
		),
		requester:          params.Requester,
		sourcesProvider:    params.SourcesProvider,
		workers:            params.Workers,
		checkInterval:      params.CheckInterval,
		responseProcessors: params.ResponseProcessors,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.checkInterval)
	for {
		select {
		case <-ctx.Done():
			r.logger.InfoContext(ctx, "context done - stopping healthcheck runner")
			return nil
		case <-ticker.C:
			r.logger.InfoContext(ctx, "starting healthcheck runs")
			start := time.Now()
			if err := r.check(ctx); err != nil {
				r.logger.With(
					slog.Any("error", err),
				).ErrorContext(ctx, "failed to check sources")
			}
			completedIn := time.Since(start)
			r.logger.With(
				slog.Duration("completedIn", completedIn),
				slog.Int("workers", int(r.workers)),
			).InfoContext(ctx, "completed healthcheck runs")
		}
	}
}

func (r *Runner) check(ctx context.Context) error {
	sources, err := r.sourcesProvider.Provide(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to provide sources")
	}
	if r.workers > 1 {
		return r.runParallel(ctx, sources)
	}
	return r.run(ctx, sources)
}

func (r *Runner) runParallel(ctx context.Context, sources types.Sources) error {
	sourcesBatches := r.generateBatches(sources)
	errGroup, groupCtx := errgroup.WithContext(ctx)
	for _, batch := range sourcesBatches {
		// this if fixed in 1.22 release
		b := batch
		errGroup.Go(func() error {
			return r.run(groupCtx, b)
		})
	}
	if err := errGroup.Wait(); err != nil {
		return errors.Wrap(err, "failed to run healthcheck in parallel")
	}
	return nil
}

func (r *Runner) run(ctx context.Context, sources types.Sources) error {
	responsesChan := r.sendHealthChecks(ctx, sources)
	r.consumeResponses(ctx, responsesChan)
	return nil
}

func (r *Runner) sendHealthChecks(ctx context.Context, sources []*url.URL) <-chan types.HealthResponse {
	responsesChan := make(chan types.HealthResponse, len(sources))
	go func() {
		for _, source := range sources {
			select {
			case <-ctx.Done():
				close(responsesChan)
				return
			default:
			}
			resp, err := r.requester.Request(ctx, source)
			if err != nil {
				r.logger.With(
					slog.Any("error", err),
					slog.String("source", source.String()),
				).ErrorContext(ctx, "failed to request source health")
				continue
			}
			responsesChan <- resp
		}
		close(responsesChan)
	}()
	return responsesChan
}

func (r *Runner) consumeResponses(ctx context.Context, responsesChan <-chan types.HealthResponse) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-responsesChan:
			if !ok {
				return
			}
			for _, processor := range r.responseProcessors {
				if err := processor.Process(ctx, resp); err != nil {
					r.logger.With(
						slog.Any("error", err),
						slog.String("source", resp.Source),
						slog.Int("statusCode", resp.StatusCode),
						slog.String("status", resp.Status),
					).ErrorContext(ctx, "failed to process response")
				}
			}
		}
	}
}

func (r *Runner) generateBatches(sources types.Sources) []types.Sources {
	batchSize := len(sources) / int(r.workers)
	if len(sources)%int(r.workers) != 0 {
		batchSize++
	}

	var batches []types.Sources
	for i := 0; i < len(sources); i += batchSize {
		end := i + batchSize
		if end > len(sources) {
			end = len(sources)
		}
		batch := sources[i:end]
		batches = append(batches, batch)
	}
	return batches
}
