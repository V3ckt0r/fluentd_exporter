FROM golang:1.8

RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY ./fluentd_exporter.go .

RUN go-wrapper download && \
    go-wrapper install

CMD ["go-wrapper","run"]
