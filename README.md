# Fluentd Exporter for Prometheus [![Build Status](https://travis-ci.org/V3ckt0r/fluentd_exporter.svg?branch=master)](https://travis-ci.org/V3ckt0r/fluentd_exporter)

Fluentd exporter uses the fluentd monitoring agent api. Documentation on setting this up can be found [here](https://docs.fluentd.org/v0.12/articles/monitoring)

Help on flags:
```
  -Telementry.endpoint string
    	Path under which to expose metric. (default "/metrics")
  -insecure
    	Ignore server certificate if using https, Default: false.
  -log.format value
    	Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] (default "info")
  -scrape_uri string
    	URI to fluentd metrics (default "http://localhost:24220/api/plugins.json")
  -telementry.address string
    	Address on which to expose metrics. (default ":9309")
  -version
    	Print version information.
```

## Collectors
The exporter collects the following metrics:

Fluentd metrics:
```
# HELP fluentd_buffer_queue_length Buffered queue length
# TYPE fluentd_buffer_queue_length counter
# HELP fluentd_buffer_total_queued_size size of the total queued
# TYPE fluentd_buffer_total_queued_size counter
# HELP fluentd_retry_count fluentd retry count
# TYPE fluentd_retry_count counter
# HELP fluentd_up Could fluentd be reached
# TYPE fluentd_up gauge
```

Request metrics:

```
# HELP http_request_duration_microseconds The HTTP request latencies in microseconds.
# TYPE http_request_duration_microseconds summary
# HELP http_request_size_bytes The HTTP request sizes in bytes.
# TYPE http_request_size_bytes summary
# HELP http_response_size_bytes The HTTP response sizes in bytes.
# TYPE http_response_size_bytes summary
```

## Building and Running
```
go build fluentd-exporter.go
./fluentd-exporter
```

## Grafana example
The following [dashbaord](https://grafana.com/dashboards/3522) can be imported into grafana
![Grafana dashboard](https://i.imgur.com/oBY6urR.png)

## Docker
It is intended that this exporter be placed on the host that Fluentd is running from. If you are using fluentd as a log driver for Docker then place this exporter on the host.
Documentation around setting up Fluentd as a logging driver for Docker can be found [here](https://docs.docker.com/engine/admin/logging/fluentd/)

## Contribute 
Feel free to open an issue or PR if you have suggestions or ideas about what to add.

## Author
[Burhan Deniz Abdi](http://www.burhan.io/)

