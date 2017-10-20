# Fluentd Exporter for Prometheus

Fluentd exporter uses the fluentd monitoring agent api. Documentation on setting this up can be found [here](https://docs.fluentd.org/v0.12/articles/monitoring)

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

## Running
```
go run fluentd-exporter.go -insecure
```

## Todo:
- add prometheus/promhttp
