package client

import (
	"strings"
	"testing"
)

const (
	basePath              = "https://test-api.configcat.com"
	invalid               = "invalid"
	testBasicAuthUsername = "08d86d63-2c48-4fbe-88a2-60e1acb40b45"
	testBasicAuthPassword = "XeWNy2BOaeRBfaZaZOp/6/t1ck9CQzxr5g6xrVhnKyE="
)

func TestClientFails(t *testing.T) {
	_, err := NewClient(basePath, invalid, invalid, "dev")
	if !strings.HasPrefix(err.Error(), "401 Unauthorized") {
		t.Errorf("Expected 401 Unauthorized. Received %s", err)
	}
}

func TestClientWorks(t *testing.T) {
	_, err := NewClient(basePath, testBasicAuthUsername, testBasicAuthPassword, "dev")
	if err != nil {
		t.Error(err)
	}
}
