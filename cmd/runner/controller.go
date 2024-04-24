package runner

import (
	"context"
	"log/slog"

	"github.com/adamkisala/health/internal/healthcheck"
)

type Controller struct {
	runner *healthcheck.Runner
	logger *slog.Logger
}

type ControllerParams struct {
	Runner *healthcheck.Runner
	Logger *slog.Logger
}

func NewController(params ControllerParams) *Controller {
	return &Controller{
		runner: params.Runner,
		logger: params.Logger.With(
			slog.String("package", "runner"),
			slog.String("struct", "Controller"),
		),
	}
}

func (c *Controller) Run(ctx context.Context) error {
	c.logger.InfoContext(ctx, "starting health check runner")
	err := c.runner.Run(ctx)
	if err != nil {
		c.logger.ErrorContext(ctx, "health check runner failed", slog.Any("error", err))
		return err
	}
	c.logger.InfoContext(ctx, "health check runner stopped")
	return nil
}
