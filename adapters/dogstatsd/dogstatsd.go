package adapters

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gliderlabs/logspout/router"
)

func init() {
	router.AdapterFactories.Register(newDogstatsdAdapter, "metrics")
}

func newDogstatsdAdapter(route *router.Route) (router.LogAdapter, error) {
	c, err := statsd.NewBuffered("127.0.0.1:8126", 1000)
	if err != nil {
		return nil, fmt.Errorf("error initializing dogstatsd client: %v", err)
	}
	return &dogstatsdAdapter{
		statsd: c,
	}, nil
}

type dogstatsdAdapter struct {
	statsd *statsd.Client
}

func (a *dogstatsdAdapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		a.inc(m)
	}
}

func (a *dogstatsdAdapter) inc(m *router.Message) {
	tags := []string{
		fmt.Sprintf("image_name:%s", m.Container.Image),
		fmt.Sprintf("container_name:%s", m.Container.Name),
		fmt.Sprintf("container_id:%s", m.Container.ID),
	}
	for k, v := range m.Container.Config.Labels {
		tags = append(tags, fmt.Sprintf("%s:%s", k, v))
	}
	a.statsd.Count("logspout.message", 1, tags, 1.0)
}
