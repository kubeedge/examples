package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIsSuccessStatus(t *testing.T) {
	for _, test := range []struct {
		statusCode int
		expected   bool
	}{
		{statusCode: 199, expected: false},
		{statusCode: 200, expected: true},
		{statusCode: 204, expected: true},
		{statusCode: 300, expected: false},
	} {
		if actual := isSuccessStatus(test.statusCode); actual != test.expected {
			t.Errorf("isSuccessStatus(%d) = %v, want %v", test.statusCode, actual, test.expected)
		}
	}
}

func TestHandlersRejectNonSuccessUpstreamStatus(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "upstream unavailable", http.StatusServiceUnavailable)
	}))
	defer upstream.Close()

	previousBaseURL := edgeAPIBaseURL
	edgeAPIBaseURL = upstream.URL
	defer func() {
		edgeAPIBaseURL = previousBaseURL
	}()

	gin.SetMode(gin.TestMode)
	tests := []struct {
		name    string
		handler gin.HandlerFunc
	}{
		{name: "get counter", handler: GetCounter},
		{name: "get status", handler: GetStatus},
		{name: "set status", handler: SetStatus},
		{name: "reset count", handler: ResetCount},
		{name: "set count", handler: SetCount},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Request = httptest.NewRequest(http.MethodGet, "/?status=on&count=1", nil)

			test.handler(context)

			if recorder.Code != http.StatusBadGateway {
				t.Fatalf("handler returned status %d, want %d", recorder.Code, http.StatusBadGateway)
			}
			if !strings.Contains(recorder.Body.String(), "503 Service Unavailable") {
				t.Fatalf("handler response %q does not expose the upstream status", recorder.Body.String())
			}
		})
	}
}
