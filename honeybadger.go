package main

import (
	"os"

	honeybadger "github.com/honeybadger-io/honeybadger-go"
)

func init() {
	honeybadger.Configure(honeybadger.Configuration{
		APIKey: os.Getenv("HONEYBADGER_API_KEY"),
		Env:    os.Getenv("HONEYBADGER_ENVIRONMENT"),
	})
}
