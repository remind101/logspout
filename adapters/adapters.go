package adapters

import (
	"sync"
	"time"

	"github.com/gliderlabs/logspout/router"
	"github.com/remind101/metrics"
)

func init() {
	router.AdapterFactories.Register(newMetricsAdapter, "metrics")
}

func newMetricsAdapter(route *router.Route) (router.LogAdapter, error) {
	a := &metricsAdapter{}
	go a.start()
	return a, nil
}

type metricsAdapter struct {
	sync.Mutex
	count int32 // Current count of messages received.
}

func (a *metricsAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		a.add(m)
	}
}

func (a *metricsAdapter) add(m *router.Message) {
	a.Lock()
	defer a.Unlock()
	a.count += 1
}

func (a *metricsAdapter) flush() {
	a.Lock()
	defer a.Unlock()
	metrics.Count("logspout.messages", a.count)
	a.count = 0
}

func (a *metricsAdapter) start() {
	for {
		<-time.After(10 * time.Second)
		a.flush()
	}
}
