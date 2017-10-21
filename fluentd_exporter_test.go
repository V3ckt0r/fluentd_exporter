package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	metricCount = 4
)

var (
	fluentdJson = Message{1, 4, 456, 0}
)

type Message struct {
	up                    int
	bufferQueueLength     int
	bufferTotalQueuedSize float64
	retryCount            int
}

func checkFluentdStatus(t *testing.T, status []byte, metricCount int) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(status))
	})
	server := httptest.NewServer(handler)

	e := NewExporter(server.URL)
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	m := <-ch
	if m == nil {
		t.Error("expected metric but got nil")
	}

	if <-ch != nil {
		t.Error("expected closed channel")
	}
}

func TestFluentdStatus(t *testing.T) {
	//marshal json data
	data, err := json.Marshal(fluentdJson)
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}

	checkFluentdStatus(t, data, metricCount)
}
