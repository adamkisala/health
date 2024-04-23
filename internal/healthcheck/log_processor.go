package healthcheck

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
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
	switch response.StatusCode {
	case http.StatusOK:
		if response.ResponseTime > p.acceptableResponseTime {
			loggerWithFields.WarnContext(ctx, "response time is too long")
		}
		for _, component := range response.Components {
			if !strings.EqualFold(component.Status, "ok") {
				loggerWithFields.ErrorContext(ctx, "component is not ok", slog.String("component", component.Name))
			}
		}
		loggerWithFields.InfoContext(ctx, "health check passed")
	case http.StatusNotFound:
		loggerWithFields.WarnContext(ctx, "health check failed - not found")
	default:
		loggerWithFields.ErrorContext(ctx, "health check failed - unexpected status code")
	}
}
