package healthcheck

import (
	"context"
	"log/slog"

	"github.com/adamkisala/health/types"
)

type EmailAlerterProcessor struct {
	logger *slog.Logger
}

type EmailAlerterProcessorParams struct {
	Logger *slog.Logger
}

func NewEmailAlerterProcessor(params EmailAlerterProcessorParams) *EmailAlerterProcessor {
	return &EmailAlerterProcessor{
		logger: params.Logger.With(
			slog.String("package", "healthchecker"),
			slog.String("struct", "EmailAlerterProcessor"),
		),
	}
}

// Process would send an email alert if the health status is not OK
func (p EmailAlerterProcessor) Process(ctx context.Context, response types.HealthResponse) error {
	loggerWithFields := p.logger.With(
		slog.String("source", response.Source),
		slog.String("responseTime", response.ResponseTime.String()),
		slog.Time("timeStamp", response.SentAt),
		slog.String("status", response.Status),
		slog.Int("statusCode", response.StatusCode),
	)
	healthStatus := response.HealthStatus()
	switch healthStatus {
	case types.HealthStatusPartiallyAvailable, types.HealthStatusNotFound, types.HealthStatusDown:
		// todo implement email alerting logic
		loggerWithFields.With(
			slog.String("healthStatus", healthStatus.String()),
			slog.String("source", response.Source)).InfoContext(ctx, "health check failed - email sent")
	default:
	}
	return nil
}
