package sources

import (
	"context"
	"io"
	"io/fs"
	"log/slog"
	"net/url"

	"github.com/adamkisala/weaviate-health/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ConfigSources represents the configuration stored
// within the config yaml file
type ConfigSources struct {
	Sources []string `yaml:"sources"`
}

// Config loads sources for health checker from
// service configuration
type Config struct {
	sourcesStore         fs.FS
	logger               *slog.Logger
	sourcesStoreFilePath string
}

// ConfigParams holds the parameters for the Config source
type ConfigParams struct {
	SourcesStore         fs.FS
	SourcesStoreFilePath string
	Logger               *slog.Logger
}

// NewConfig creates a new Config source
func NewConfig(params ConfigParams) *Config {
	if params.Logger == nil {
		params.Logger = slog.Default()
	}
	return &Config{
		sourcesStore:         params.SourcesStore,
		sourcesStoreFilePath: params.SourcesStoreFilePath,
		logger: params.Logger.With(
			slog.String("package", "sources"),
			slog.String("struct", "Config")),
	}
}

// Load loads the sources from the Config source
func (c *Config) Load(ctx context.Context) (types.Sources, error) {
	file, err := c.sourcesStore.Open(c.sourcesStoreFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open config file")
	}
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	var configSources ConfigSources
	if err := yaml.Unmarshal(fileContents, &configSources); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config file")
	}
	srcs := make(types.Sources, 0, len(configSources.Sources))
	for _, configSource := range configSources.Sources {
		srcURL, err := url.Parse(configSource)
		if err != nil {
			// we still want to load remaining sources, so log the error and continue
			c.logger.With(
				slog.Any("error", err),
				slog.String("configSource", configSource)).
				ErrorContext(ctx, "failed to parse source URL")
			continue
		}
		srcs = append(srcs, srcURL)
	}
	return srcs, nil
}
