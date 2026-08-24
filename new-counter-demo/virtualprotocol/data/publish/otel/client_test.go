package otel

import (
	"encoding/json"
	"testing"
)

func TestNewConfigValidatesEndpointURL(t *testing.T) {
	if _, err := NewConfig(json.RawMessage(`{"endpointURL":"localhost:4318"}`)); err == nil {
		t.Fatal("NewConfig accepted an endpoint without an HTTP URL host")
	}
	if _, err := NewConfig(json.RawMessage(`{"endpointURL":"http://localhost:4318/v1/metrics"}`)); err != nil {
		t.Fatalf("NewConfig rejected a valid endpoint: %v", err)
	}
}
