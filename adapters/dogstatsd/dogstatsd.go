package dogstatsd

import (
	"fmt"
	"strings"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gliderlabs/logspout/router"
)

func init() {
	router.AdapterFactories.Register(newDogstatsdAdapter, "dogstatsd")
}

func newDogstatsdAdapter(route *router.Route) (router.LogAdapter, error) {
	c, err := statsd.New(route.Address)
	if err != nil {
		return nil, fmt.Errorf("error initializing dogstatsd client: %v", err)
	}
	labels := labelsFromString(route.Options["labels"])
	return &dogstatsdAdapter{
		statsd: c,
		labels: labels,
	}, nil
}

type dogstatsdAdapter struct {
	statsd *statsd.Client

	// List of labels to include as tags on the metric.
	labels map[string]bool
}

func (a *dogstatsdAdapter) Stream(logstream chan *router.Message) {
	var logspout_dogstatsd_enabled = os.Getenv("LOGSPOUT_DOGSTATSD_ENABLED")
	if logspout_dogstatsd_enabled != "false" {
		for m := range logstream {
			a.inc(m)
		}
	}
}

func (a *dogstatsdAdapter) inc(m *router.Message) {
	tags := []string{
		fmt.Sprintf("image_name:%s", m.Container.Config.Image),
	}
	for name := range a.labels {
		if v, ok := m.Container.Config.Labels[name]; ok {
			tags = append(tags, fmt.Sprintf("%s:%s", name, v))
		}
	}
	a.statsd.Count("logspout.message", 1, tags, 1.0)
	a.statsd.Histogram("logspout.message.size", float64(len(m.Data)), tags, 1.0)
}

// Splits a comma separated list of labels into a map[string]bool.
func labelsFromString(s string) map[string]bool {
	if s == "" {
		return nil
	}

	labels := make(map[string]bool)
	for _, name := range strings.Split(s, ",") {
		labels[name] = true
	}
	return labels
}
