FROM golang:1.10.1-alpine

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY ./fluentd_exporter .

CMD ["fluentd_exporter"]
