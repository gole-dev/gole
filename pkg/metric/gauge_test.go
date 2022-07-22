package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewGaugeVec(t *testing.T) {
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_server",
		Subsystem: "requests",
		Name:      "duration",
		Help:      "rpc server requests duration(ms).",
	})
	gaugeVecNil := NewGaugeVec(nil)
	assert.NotNil(t, gaugeVec)
	assert.Nil(t, gaugeVecNil)
}

func TestGaugeInc(t *testing.T) {
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_client2",
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"path"},
	})
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Inc("/users")
	gv.Inc("/users")
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(2), r)
}

func TestGaugeAdd(t *testing.T) {
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "rpc_client",
		Subsystem: "request",
		Name:      "duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"path"},
	})
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Add(-10, "/liveroom")
	gv.Add(30, "/liveroom")
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(20), r)
}

func TestGaugeSet(t *testing.T) {
	gaugeVec := NewGaugeVec(&GaugeVecOpts{
		Namespace: "http_client",
		Subsystem: "request",
		Name:      "duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"path"},
	})
	gv, _ := gaugeVec.(*promGaugeVec)
	gv.Set(666, "/users")
	r := testutil.ToFloat64(gv.gauge)
	assert.Equal(t, float64(666), r)
}
