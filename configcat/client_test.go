package configcat

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

func TestClient_Fails(t *testing.T) {
	_, err := NewClient(basePath, invalid, invalid)
	if !strings.HasPrefix(err.Error(), "401 Unauthorized") {
		t.Errorf("Expected 401 Unauthorized. Received %s", err)
	}
}

func TestClient_Works(t *testing.T) {
	_, err := NewClient(basePath, testBasicAuthUsername, testBasicAuthPassword)
	if err != nil {
		t.Error(err)
	}
}
