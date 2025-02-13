package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	processedTickGauge prometheus.Gauge
	eventTickGauge     prometheus.Gauge
	liveTickGauge      prometheus.Gauge
	liveEpochGauge     prometheus.Gauge
}

func NewMetrics() *Metrics {
	m := Metrics{
		processedTickGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "qubic_transfers_processed_tick",
			Help: "The latest fully processed tick",
		}),
		eventTickGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "qubic_transfers_event_tick",
			Help: "The latest known event tick",
		}),
		liveTickGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "qubic_transfers_live_tick",
			Help: "The latest known live tick",
		}),
	}
	return &m
}

func (metrics *Metrics) SetLatestProcessedTick(tick uint32) {
	metrics.processedTickGauge.Set(float64(tick))
}

func (metrics *Metrics) SetLatestAvailableEventTick(tick uint32) {
	metrics.eventTickGauge.Set(float64(tick))
}

func (metrics *Metrics) SetLatestAvailableLiveTick(tick uint32) {
	metrics.liveTickGauge.Set(float64(tick))
}
