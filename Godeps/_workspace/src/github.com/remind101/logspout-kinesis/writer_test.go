package kinesis

import (
	"testing"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/gliderlabs/logspout/router"
)

var m = &router.Message{
	Data: "hello",
}

type fakeFlusher struct {
	inputs        chan kinesis.PutRecordsInput
	dropInputFunc func(kinesis.PutRecordsInput)
	flushFunc     func()
	flushed       chan struct{}
}

func (f *fakeFlusher) start() {
	f.flushInputs()
}

func (f *fakeFlusher) flush(input kinesis.PutRecordsInput) {
	select {
	case f.inputs <- input:
	default:
		f.dropInputFunc(input)
	}
}

func (f *fakeFlusher) flushInputs() {
	if f.flushFunc == nil {
		close(f.flushed)
	} else {
		f.flushFunc()
	}
}

var testLimits = limits{
	putRecords:     2,
	putRecordsSize: PutRecordsSizeLimit,
	recordSize:     RecordSizeLimit,
}

func TestWriter_Flush(t *testing.T) {
	tmpl, _ := template.New("").Parse("abc")
	b := newBuffer(tmpl, "abc")
	b.limits = &testLimits

	f := &fakeFlusher{
		inputs:  make(chan kinesis.PutRecordsInput, 10),
		flushed: make(chan struct{}),
	}

	w := newWriter(b, f)
	w.ticker = nil

	w.start()

	w.write(m)
	w.write(m)
	w.write(m)

	select {
	case <-f.flushed:
	case <-time.After(1 * time.Second):
		t.Fatal("Expected flush to be called")
	}
}

func TestWriter_PeriodicFlush(t *testing.T) {
	tmpl, _ := template.New("").Parse("abc")
	b := newBuffer(tmpl, "abc")
	b.limits = &testLimits

	f := &fakeFlusher{
		inputs:  make(chan kinesis.PutRecordsInput, 10),
		flushed: make(chan struct{}),
	}

	w := newWriter(b, f)

	ticker := make(chan time.Time)
	w.ticker = ticker

	w.start()
	w.write(m)

	select {
	case ticker <- time.Now():
	default:
		t.Fatal("Couldn't send on stream.ticker channel")
	}

	select {
	case <-f.flushed:
	case <-time.After(time.Second):
		t.Fatal("Expected flush to be called")
	}
}
