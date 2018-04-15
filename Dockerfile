FROM golang:1.10.1-alpine

COPY ./fluentd_exporter /bin/fluentd_exporter

ENTRYPOINT ["/bin/fluentd_exporter"]
