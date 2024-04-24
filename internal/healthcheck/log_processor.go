package healthcheck

import (
	"context"
	"log/slog"
	"time"

	"github.com/adamkisala/weaviate-health/types"
)

type LogProcessor struct {
	logger                 *slog.Logger
	acceptableResponseTime time.Duration
}

type LogProcessorParams struct {
	Logger                 *slog.Logger
	AcceptableResponseTime time.Duration
}

func NewLogProcessor(params LogProcessorParams) *LogProcessor {
	return &LogProcessor{
		logger: params.Logger.With(
			slog.String("package", "healthchecker"),
			slog.String("struct", "LogProcessor"),
		),
		acceptableResponseTime: params.AcceptableResponseTime,
	}
}

func (p LogProcessor) Process(ctx context.Context, response types.HealthResponse) error {
	loggerWithFields := p.logger.With(
		slog.String("source", response.Source),
		slog.String("responseTime", response.ResponseTime.String()),
		slog.Time("timeStamp", response.TimeStamp),
		slog.String("status", response.Status),
		slog.Int("statusCode", response.StatusCode),
	)
	p.logResponse(ctx, response, loggerWithFields)
	return nil
}

func (p LogProcessor) logResponse(ctx context.Context, response types.HealthResponse, loggerWithFields *slog.Logger) {
	if response.ResponseTime > p.acceptableResponseTime {
		loggerWithFields.WarnContext(ctx, "response time is too long")
	}
	healthStatus := response.HealthStatus()
	switch healthStatus {
	case types.HealthStatusOK:
		loggerWithFields.With(
			slog.String("healthStatus", healthStatus.String())).
			InfoContext(ctx, "health check passed")
	case types.HealthStatusNotFound:
		loggerWithFields.With(
			slog.String("healthStatus", healthStatus.String())).
			WarnContext(ctx, "health check failed")
	case types.HealthStatusPartiallyAvailable:
		loggerWithFields.With(
			slog.String("healthStatus", healthStatus.String())).
			ErrorContext(ctx, "health check failed")
	default:
		loggerWithFields.With(
			slog.String("healthStatus", healthStatus.String())).
			ErrorContext(ctx, "health check failed - unexpected status code")
	}
}
