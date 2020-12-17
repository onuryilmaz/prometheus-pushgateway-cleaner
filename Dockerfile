FROM golang:1.14 as builder
ADD . /go/src/github.com/onuryilmaz/prometheus-pushgateway-cleaner/
WORKDIR /go/src/github.com/onuryilmaz/prometheus-pushgateway-cleaner/cmd
RUN GOOS=linux go build -o prometheus-pushgateway-cleaner

FROM alpine
COPY --from=builder /go/src/github.com/onuryilmaz/prometheus-pushgateway-cleaner/cmd/prometheus-pushgateway-cleaner /prometheus-pushgateway-cleaner
ENTRYPOINT ["./prometheus-pushgateway-cleaner"]