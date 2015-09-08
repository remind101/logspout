package main

import (
	_ "github.com/gliderlabs/logspout/adapters/raw"
	_ "github.com/gliderlabs/logspout/httpstream"
	_ "github.com/gliderlabs/logspout/routesapi"
	_ "github.com/gliderlabs/logspout/transports/tcp"
	_ "github.com/remind101/logspout-kinesis"
	_ "github.com/remind101/logspout/adapters/syslog"
	_ "github.com/remind101/logspout/transports/udp"
)
