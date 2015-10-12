package main

import (
	"os"

	honeybadger "github.com/honeybadger-io/honeybadger-go"
	"github.com/remind101/logspout-kinesis"
)

func init() {
	if apiKey := os.Getenv("HONEYBADGER_API_KEY"); apiKey != "" {
		honeybadger.Configure(honeybadger.Configuration{
			APIKey: apiKey,
			Env:    os.Getenv("HONEYBADGER_ENVIRONMENT"),
		})

		kinesis.ErrorHandler = func(err error) {
			if err != nil {
				if _, ok := err.(*kinesis.StreamNotReadyError); !ok {
					honeybadger.DefaultClient.Notify(err)
				}
			}
		}
	}
}
