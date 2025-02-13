package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-transfers/metrics"
	"net/http"
	"time"
)

type MetricsServer struct {
	address        string
	metricsService *metrics.Metrics
}

func NewMetricsServer(address string, metricsService *metrics.Metrics) *MetricsServer {
	server := &MetricsServer{
		address:        address,
		metricsService: metricsService,
	}
	return server
}

func (s *MetricsServer) Start() {

	go func() {
		serverMux := http.NewServeMux()
		serverMux.Handle("/metrics", promhttp.Handler()) // FIXME

		var server = &http.Server{
			Addr:              s.address,
			Handler:           serverMux,
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 15 * time.Second,
			WriteTimeout:      15 * time.Second,
		}

		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}
