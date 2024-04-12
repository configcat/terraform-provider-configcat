package client

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRetry(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.Header().Add("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
		}))
	defer ts.Close()

	client := &http.Client{
		Transport: Retry(http.DefaultTransport, 2),
	}

	resp, err := client.Post(ts.URL, "text/plain", strings.NewReader("body"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	assert.Equal(t, 3, attempts)
}

func TestRetry_Eventually_Ok(t *testing.T) {
	attempts := 0
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			if attempts <= 2 {
				w.Header().Add("Retry-After", "1")
				w.WriteHeader(http.StatusTooManyRequests)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}))
	defer ts.Close()

	client := &http.Client{
		Transport: Retry(http.DefaultTransport, 10),
	}

	resp, err := client.Post(ts.URL, "text/plain", strings.NewReader("body"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, attempts)
}
