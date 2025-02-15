package metrics

import (
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

var meters = NewMetrics()

func TestEventService_SetLatestAvailableEventTick(t *testing.T) {
	meters.SetLatestEventTick(42)
	assert.Equal(t, float64(42), testutil.ToFloat64(meters.eventTickGauge))
}

func TestEventService_SetLatestProcessedTick(t *testing.T) {
	meters.SetLatestProcessedTick(43)
	assert.Equal(t, float64(43), testutil.ToFloat64(meters.processedTickGauge))
}

func TestEventService_SetLatestAvailableLiveTick(t *testing.T) {
	meters.SetLatestLiveTick(44)
	assert.Equal(t, float64(44), testutil.ToFloat64(meters.liveTickGauge))
}
