package healthcheck

import (
	"net/http"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
)

var transientErrorCodes = map[int]bool{
	http.StatusRequestTimeout:     true,
	http.StatusTooManyRequests:    true,
	http.StatusBadGateway:         true,
	http.StatusServiceUnavailable: true,
	http.StatusGatewayTimeout:     true,
}

type RetryableClient struct {
	maxRetries uint64
	retryWait  time.Duration
	client     *http.Client
}

type RetryableClientParams struct {
	MaxRetries uint64
	RetryWait  time.Duration
	Client     *http.Client
}

func NewRetryableClient(params RetryableClientParams) *RetryableClient {
	return &RetryableClient{
		maxRetries: params.MaxRetries,
		retryWait:  params.RetryWait,
		client:     params.Client,
	}
}

func (r *RetryableClient) Do(req *http.Request) (*http.Response, error) {
	bOff := backoff.WithContext(
		backoff.WithMaxRetries(
			backoff.NewConstantBackOff(r.retryWait), r.maxRetries),
		req.Context())

	var resp *http.Response
	operation := func() error {
		var err error
		resp, err = r.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "failed to perform request")
		}
		if ok := transientErrorCodes[resp.StatusCode]; ok {
			return errors.Errorf("transient error: %d", resp.StatusCode)
		}
		return nil
	}
	if err := backoff.Retry(operation, bOff); err != nil {
		return nil, errors.Wrap(err, "failed to perform request with backoff")
	}
	return nil, nil
}
