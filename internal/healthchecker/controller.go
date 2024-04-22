package healthchecker

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
	Provide(ctx context.Context) ([]*url.URL, error)
}

type Runner struct {
	logger            *slog.Logger
	requester         Requester
	responseProcessor ResponseProcessor
	sourcesProvider   SourcesProvider
	workers           int64
	checkInterval     time.Duration
}

type RunnerParams struct {
	Logger            *slog.Logger
	Requester         Requester
	ResponseProcessor ResponseProcessor
	SourcesProvider   SourcesProvider
	Workers           int64
	CheckInterval     time.Duration
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
		requester:         params.Requester,
		sourcesProvider:   params.SourcesProvider,
		workers:           params.Workers,
		checkInterval:     params.CheckInterval,
		responseProcessor: params.ResponseProcessor,
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
			).InfoContext(ctx, "completed healthcheck run")
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

func (r *Runner) runParallel(ctx context.Context, sources []*url.URL) error {
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

func (r *Runner) run(ctx context.Context, sources []*url.URL) error {
	sourcesChan := make(chan *url.URL, len(sources))
	for _, source := range sources {
		sourcesChan <- source
	}
	responsesChan := make(chan types.HealthResponse, len(sources))

	errGroup, groupCtx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		defer close(sourcesChan)
		r.sendHealthChecks(groupCtx, sourcesChan, responsesChan)
		return nil
	})
	errGroup.Go(func() error {
		defer close(responsesChan)
		r.consumeResponses(groupCtx, responsesChan)
		return nil
	})
	if err := errGroup.Wait(); err != nil {
		return errors.Wrap(err, "failed to run healthcheck")
	}
	return nil
}

func (r *Runner) sendHealthChecks(ctx context.Context, sourcesChan <-chan *url.URL, responsesChan chan<- types.HealthResponse) {
	for source := range sourcesChan {
		select {
		case <-ctx.Done():
			return
		default:
		}
		resp, err := r.requester.Request(ctx, source)
		if err != nil {
			r.logger.With(
				slog.Any("error", err),
				slog.String("source", source.String()),
			).ErrorContext(ctx, "failed to request source health")
		}
		responsesChan <- resp
	}
}

func (r *Runner) consumeResponses(ctx context.Context, responsesChan chan types.HealthResponse) {
	for resp := range responsesChan {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if err := r.responseProcessor.Process(ctx, resp); err != nil {
			r.logger.With(
				slog.Any("error", err),
				slog.String("source", resp.Source),
				slog.Int("statusCode", resp.StatusCode),
				slog.String("status", resp.Status),
			).ErrorContext(ctx, "failed to process response")
		}
	}
}

func (r *Runner) generateBatches(sources []*url.URL) [][]*url.URL {
	batchSize := len(sources) / int(r.workers)
	if len(sources)%int(r.workers) != 0 {
		batchSize++
	}

	var batches [][]*url.URL
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
