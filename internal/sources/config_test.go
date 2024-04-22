package sources_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/adamkisala/weaviate-health/internal/sources"
	"github.com/adamkisala/weaviate-health/types"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	type fields struct {
		configFilePath string
		logger         *slog.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		testFileFuncCloser func(t *testing.T, content string) func() error
		want               types.Sources
		errLoad            error
	}{
		{
			name: "should return error when config file path is empty",
			fields: fields{
				configFilePath: "",
			},
			errLoad: sources.ErrorConfigFilePathRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := sources.NewConfig(sources.ConfigParams{
				ConfigFilePath: tt.fields.configFilePath,
				Logger:         tt.fields.logger,
			})
			result, err := c.Load(tt.args.ctx)
			assert.ErrorIs(t, err, tt.errLoad)
			assert.Equal(t, result, tt.want)
		})
	}
}
