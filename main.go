package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

var (
	httpRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "HTTP requests count",
	})
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Info().Msg("go-Prometheus Client app starting...")

	prometheus.MustRegister(httpRequests)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpRequests.Inc()

		clientIP := r.RemoteAddr
		log.Info().
			Str("method", r.Method).
			Str("path", "/").
			Str("client_ip", clientIP).
			Msg("Main Page request received")

		_, err := w.Write([]byte("Hello, Prometheus Client run successful with Zerolog!"))
		if err != nil {
			log.Error().Err(err).Msg("Failed to write response")
		}
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Info().Msg(" Started Client on port 8080!")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("Error" +
			" start client!")
	}

}
