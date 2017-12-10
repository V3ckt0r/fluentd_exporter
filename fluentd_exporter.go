package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"

	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"sync"
)

const (
	namespace = "fluentd" //for Prometheus metrics.
)

// declare variables for fluentd metrics
var (
	listeningAddress = flag.String("telemetry.address", ":9309", "Address on which to expose metrics.")
	metricsEndpoint  = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metric.")
	scrapeURI        = flag.String("scrape_uri", "http://localhost:24220/api/plugins.json", "URI to fluentd metrics")
	insecure         = flag.Bool("insecure", false, "Ignore server certificate if using https, Default: false.")
	showVersion      = flag.Bool("version", false, "Print version information.")
)

type Exporter struct {
	URI    string
	mutex  sync.Mutex
	client *http.Client

	up *prometheus.Desc
	bufferQueueLength *prometheus.Desc
	bufferTotalQueuedSize *prometheus.Desc
	retryCount            *prometheus.Desc
}

// NewExporter returns an initialized Exporter.

func NewExporter(uri string) *Exporter {
	return &Exporter{
		URI: uri,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Could fluentd be reached",
			nil,
			nil),
		bufferQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "buffer_queue_length"),
			"Buffered queue length",
			[]string{"pluginId", "pluginCategory"},
			nil),
		bufferTotalQueuedSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "buffer_queued_size"),
			"size of the total queued",
			[]string{"pluginId", "pluginCategory"},
			nil),
		retryCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "retry_total"),
			"fluentd retry count",
			[]string{"pluginId", "pluginCategory"},
			nil),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
			},
		},
	}
}

// Describe describes all the metrics ever exported by the fluentd exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.bufferQueueLength
	ch <- e.bufferTotalQueuedSize
	ch <- e.retryCount
}

// json data structure for fluentd
type jsonData struct {
	Plugins []struct {
		PluginID              string `json:"plugin_id"`
		PluginCategory        string `json:"plugin_category"`
		Type                  string `json:"type"`
		RetryCount            int    `json:"retry_count"`
		BufferQueueLength     int    `json:"buffer_queue_length,omitempty"`
		BufferTotalQueuedSize int    `json:"buffer_total_queued_size,omitempty"`
	} `json:"plugins"`
}

// Collect fetches the stats from configured location and delivers them
// as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	resp, err := e.client.Get(e.URI)
	if err != nil {
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		return fmt.Errorf("Error scraping fluentd: %v", err)
	}
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 1)

	// get data from body of response and check if there was a read error
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	// close connection
	resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("There was an error")
		return fmt.Errorf("Status %s (%d)", resp.Status, resp.StatusCode)
	}

	//init struct for unmarshal and check that there was no unmarshalling error
	jdata := jsonData{}
	jError := json.Unmarshal(body, &jdata)
	if jError != nil {
		log.Fatal(jError)
	}

	// organise json response and map to created metrics
	//parse through json response from fluentd
	for _, plugin := range jdata.Plugins {
		if plugin.PluginCategory == "input" {
			continue
		}
		ch <- prometheus.MustNewConstMetric(e.bufferQueueLength, prometheus.GaugeValue, float64(plugin.BufferQueueLength), plugin.PluginID, plugin.PluginCategory)
		ch <- prometheus.MustNewConstMetric(e.bufferTotalQueuedSize, prometheus.CounterValue, float64(plugin.BufferTotalQueuedSize), plugin.PluginID, plugin.PluginCategory)
		ch <- prometheus.MustNewConstMetric(e.retryCount, prometheus.CounterValue, float64(plugin.RetryCount), plugin.PluginID, plugin.PluginCategory)
	}

	return nil
}

// Collect fetches the stats from configured fluentd location and delivers them
// as Prometheus metrics.
// It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		log.Errorf("Error scraping fluentd: %s", err)
	}
	return
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("fluentd_exporter"))
		os.Exit(0)
	}

	exporter := NewExporter(*scrapeURI)

	prometheus.MustRegister(exporter)
	prometheus.MustRegister(version.NewCollector("fluentd_exporter"))

	log.Infoln("Starting fluentd_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	log.Infof("Starting Server: %s", *listeningAddress)
	http.Handle(*metricsEndpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
