package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_model/go"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMainPageHandler(t *testing.T) {
	// Registering the metric for testing
	log.Info().Msg("Registering the httpRequests metric for testing.")
	prometheus.MustRegister(httpRequests)
	defer prometheus.Unregister(httpRequests)

	// Creating a test request for the main page
	log.Info().Msg("Creating a test request for the main page handler.")
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Defining the handler for the main page, which increments the httpRequests counter
	log.Info().Msg("Defining and executing the main page handler with httpRequests counter increment.")
	mainPageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()
		_, err := w.Write([]byte("Hello, Prometheus with Zerolog!"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to write response in main page handler")
			t.Fatalf("unexpected error: %v", err)
		}
	})
	mainPageHandler.ServeHTTP(w, req)

	// Checking the response status and content
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	log.Info().Msg("Verifying the status code and response body of the main page handler.")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Hello, Prometheus with Zerolog!", string(body))

	// Checking that the httpRequests counter increased
	log.Info().Msg("Checking that the httpRequests counter increased.")
	metric := &io_prometheus_client.Metric{}
	err := httpRequests.Write(metric)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write metric data")
		t.Fatalf("unexpected error: %v", err)
	}
	assert.Equal(t, float64(1), *metric.Counter.Value)
}

func TestMetricsEndpoint(t *testing.T) {
	// Registering the httpRequests metric and incrementing it to ensure visibility in /metrics endpoint
	log.Info().Msg("Registering and incrementing the httpRequests metric for the /metrics endpoint.")
	prometheus.MustRegister(httpRequests)
	defer prometheus.Unregister(httpRequests)
	httpRequests.Inc()

	// Creating a test request for the /metrics endpoint
	log.Info().Msg("Creating a test request for the /metrics endpoint.")
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	// Executing the /metrics endpoint handler
	log.Info().Msg("Executing the /metrics endpoint handler and capturing response.")
	metricsHandler := promhttp.Handler()
	metricsHandler.ServeHTTP(w, req)

	// Checking that the /metrics response contains the http_requests_total metric
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	log.Info().Msg("Verifying that the /metrics response contains the http_requests_total metric.")
	assert.Contains(t, string(body), "http_requests_total")
}
