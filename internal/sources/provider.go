package sources

import (
	"context"
	"net/url"
	"sync"

	"github.com/pkg/errors"
)

type Loader interface {
	Load(ctx context.Context) ([]*url.URL, error)
}

type Provider struct {
	sources []*url.URL
	mux     *sync.RWMutex
	loader  Loader
}

type ProviderParams struct {
	Loader Loader
}

func NewProvider(params ProviderParams) *Provider {
	return &Provider{
		loader: params.Loader,
		mux:    &sync.RWMutex{},
	}
}

func (p *Provider) Provide(ctx context.Context) ([]*url.URL, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if len(p.sources) == 0 {
		if err := p.reload(ctx); err != nil {
			return nil, errors.Wrap(err, "failed to reload sources on Provide")
		}
	}
	return p.sources, nil
}

func (p *Provider) Reload(ctx context.Context) ([]*url.URL, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if err := p.reload(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to reload sources on Reload")
	}
	return p.sources, nil
}

func (p *Provider) reload(ctx context.Context) error {
	sources, err := p.loader.Load(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to load sources from provider")
	}
	p.sources = sources
	return nil
}
