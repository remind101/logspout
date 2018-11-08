package main

import (
	_ "github.com/gliderlabs/logspout/adapters/raw"
	_ "github.com/gliderlabs/logspout/transports/tcp"
	_ "github.com/remind101/logspout-kinesis"
	_ "github.com/remind101/logspout/adapters/datadog"
	_ "github.com/remind101/logspout/adapters/dogstatsd"
	_ "github.com/remind101/logspout/adapters/syslog"
	_ "github.com/remind101/logspout/transports/udp"
)
