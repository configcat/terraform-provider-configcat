package client

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

type retryRoundTripper struct {
	http.RoundTripper
	maxRetries int
}

func Retry(transport http.RoundTripper, maxRetries int) http.RoundTripper {
	return &retryRoundTripper{RoundTripper: transport, maxRetries: maxRetries}
}

func (rt *retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	bodyBuf, err := preReadRequestBodyAndClose(r)
	if err != nil {
		return nil, err
	}

	var attempts int
	for {
		if bodyBuf != nil {
			r.Body = io.NopCloser(bodyBuf)
		}
		resp, err := rt.RoundTripper.RoundTrip(r)
		if err != nil {
			return resp, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}
		retryAfterHeader := resp.Header.Get("Retry-After")
		retryAfter, err := strconv.ParseInt(retryAfterHeader, 10, 64)
		if err != nil {
			return resp, err
		}

		if attempts >= rt.maxRetries {
			return resp, nil
		}

		if bodyBuf != nil {
			if _, err := bodyBuf.Seek(0, io.SeekStart); err != nil {
				return resp, err
			}
		}

		if resp.Body != nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}

		select {
		case <-time.After(time.Duration(retryAfter) * time.Second):
			attempts++
		case <-r.Context().Done():
			return nil, r.Context().Err()
		}
	}
}

func preReadRequestBodyAndClose(r *http.Request) (*bytes.Reader, error) {
	var reader *bytes.Reader
	if r.Body != nil && r.Body != http.NoBody {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r.Body); err != nil {
			_ = r.Body.Close()
			return nil, err
		}
		_ = r.Body.Close()
		reader = bytes.NewReader(buf.Bytes())
	}
	return reader, nil
}
