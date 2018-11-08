package datadog

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"reflect"
	"text/template"
	"time"

	"github.com/gliderlabs/logspout/router"
)

func init() {
	router.AdapterFactories.Register(NewDatadogAdapter, "raw")
}

type DatadogAdapter struct {
	conn        net.Conn
	route       *router.Route
	serviceTmpl *template.Template
}

func NewDatadogAdapter(route *router.Route) (router.LogAdapter, error) {
	transport, found := router.AdapterTransports.Lookup(route.AdapterTransport("udp"))
	if !found {
		return nil, errors.New("bad transport: " + route.Adapter)
	}
	conn, err := transport.Dial(route.Address, route.Options)
	if err != nil {
		return nil, err
	}

	service := getopt("DATADOG_SERVICE", "{{.Container.Name}}")
	serviceTmpl, err := template.New("service").Parse(service)
	if err != nil {
		return nil, err
	}

	return &DatadogAdapter{
		route:       route,
		conn:        conn,
		serviceTmpl: serviceTmpl,
	}, nil
}

func (a *DatadogAdapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		buf := new(bytes.Buffer)
		err := a.serviceTmpl.Execute(buf, message)
		if err != nil {
			log.Println("datadog:", err)
		}
		service := buf.String()

		raw, err := json.Marshal(&Message{
			Timestamp: &message.Time,
			Service:   &service,
			Source:    &message.Source,
			Message:   &message.Data,
			Container: &Container{
				ID:     &message.Container.ID,
				Labels: &message.Container.Config.Labels,
			},
		})
		_, err = a.conn.Write(raw)
		if err != nil {
			log.Println("datadog:", err)
			if reflect.TypeOf(a.conn).String() != "*net.UDPConn" {
				return
			}
		}
	}
}

type Container struct {
	ID     *string            `json:"id"`
	Labels *map[string]string `json:"labels"`
}

// See the documention on "Reserved Attributes" for JSON logs
//
// https://docs.datadoghq.com/logs/processing/pipelines/#reserved-attribute-pipeline
type Message struct {
	Timestamp *time.Time `json:"timestamp"`
	Source    *string    `json:"source"`
	Service   *string    `json:"service"`
	Message   *string    `json:"message"`
	Container *Container `json:"container"`
}

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}
