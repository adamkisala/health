package sources_test

import (
	"context"
	"io/fs"
	"log/slog"
	"os"
	"testing"

	"github.com/adamkisala/health/internal/sources"
	"github.com/adamkisala/health/tests"
	"github.com/adamkisala/health/types"
	"github.com/liamg/memoryfs"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	type fields struct {
		configFilePath string
		sourcesStore   fs.FS
		logger         *slog.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Sources
		errLoad error
	}{
		{
			name: "should return error when failed to read config file",
			fields: fields{
				configFilePath: "config.yaml",
				sourcesStore:   memoryfs.New(),
			},
			errLoad: os.ErrNotExist,
		},
		{
			name: "should continue scanning sources when some are not valid",
			fields: fields{
				configFilePath: "invalid-urls-config.json",
				sourcesStore:   createTestFile(t, "invalid-urls-config.json", `sources: ["://invalid-url", "http://valid-url"]`),
			},
			want: types.Sources{
				tests.MustParseURL(t, "http://valid-url"),
			},
		},
		{
			name: "should return all valid sources when all are valid",
			fields: fields{
				configFilePath: "invalid-urls-config.json",
				sourcesStore:   createTestFile(t, "invalid-urls-config.json", `sources: ["https://valid-url.com/health", "http://valid-url"]`),
			},
			want: types.Sources{
				tests.MustParseURL(t, "https://valid-url.com/health"),
				tests.MustParseURL(t, "http://valid-url"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := sources.NewConfig(sources.ConfigParams{
				SourcesStoreFilePath: tt.fields.configFilePath,
				Logger:               tt.fields.logger,
				SourcesStore:         tt.fields.sourcesStore,
			})
			result, err := c.Load(tt.args.ctx)
			assert.ErrorIs(t, err, tt.errLoad)
			assert.Equal(t, tt.want, result)
		})
	}
}

func createTestFile(t *testing.T, s string, content string) fs.FS {
	memory := memoryfs.New()
	if err := memory.WriteFile(s, []byte(content), fs.ModePerm); err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}
	return memory
}
