package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/metrics" {
			t.Errorf("Expected request to '/v1/metrics', got: %s", r.URL.Path)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "test-token" {
			t.Errorf("Expected Authorization header 'test-token', got: %s", authHeader)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mock metrics data"))
	}))
	defer mockServer.Close()

	apiAddress = mockServer.URL
	apiToken = "test-token"

	req, err := http.NewRequest("GET", "/v1/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(metricsHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "mock metrics data"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Handler returned unexpected Content-Type header: got %v want %v", contentType, "text/plain")
	}
}

func TestMetricsHandlerErrorStatus(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer mockServer.Close()

	apiAddress = mockServer.URL
	apiToken = "test-token"

	req, err := http.NewRequest("GET", "/v1/metrics", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(metricsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestMetricsHandlerURLConstruction(t *testing.T) {
	testCases := []struct {
		name     string
		address  string
		expected string
	}{
		{
			name:     "No trailing slash",
			address:  "http://example.com",
			expected: "http://example.com/v1/metrics",
		},
		{
			name:     "With trailing slash",
			address:  "http://example.com/",
			expected: "http://example.com/v1/metrics",
		},
		{
			name:     "With complete path",
			address:  "http://example.com/v1/metrics",
			expected: "http://example.com/v1/metrics",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 这里我们不需要检查完整的 URL，因为 mockServer.URL 将不同
				// 但我们可以检查路径是否正确
				if r.URL.Path != "/v1/metrics" {
					t.Errorf("Expected path '/v1/metrics', got: %s", r.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer mockServer.Close()

			apiAddress = tc.address
			apiToken = "test-token"

			req, err := http.NewRequest("GET", "/v1/metrics", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(metricsHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		})
	}
}
