FROM golang:1.14-alpine as builder
ADD . /go/src/github.com/onuryilmaz/prometheus-pushgateway-cleaner/
WORKDIR /go/src/github.com/onuryilmaz/prometheus-pushgateway-cleaner/cmd
RUN GOOS=linux go build -o prometheus-pushgateway-cleaner

FROM scratch
COPY --from=builder /go/src/github.com/onuryilmaz/prometheus-pushgateway-cleaner/cmd/prometheus-pushgateway-cleaner /usr/bin/prometheus-pushgateway-cleaner
ENTRYPOINT ["prometheus-pushgateway-cleaner"]