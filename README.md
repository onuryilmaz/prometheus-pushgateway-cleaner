# Prometheus Pushgateway Cleaner
# :fire: :hourglass_flowing_sand: :recycle:

It is a tool to delete old metrics from [Prometheus Pushgateway](https://github.com/prometheus/pushgateway).

## Why?

From [Pushgateway Non-Goals](https://github.com/prometheus/pushgateway/blob/master/README.md#non-goals):

> A while ago, we
[decided to not implement a “timeout” or TTL for pushed metrics](https://github.com/prometheus/pushgateway/issues/19)
because almost all proposed use cases turned out to be anti-patterns we
strongly discourage.

## How?
```
$ docker run -it onuryilmaz/prometheus-pushgateway-cleaner:latest --help

Usage of prometheus-pushgateway-cleaner:
  -address string
        Address of Prometheus Pushgateway
  -debug
        Set debug log level, the default is false.
  -dry-run
        Dry run and do not delete the metrics, the default is false
  -ttl duration
        TTL for clearing the expired metrics, the default is 24 hours. (default 24h0m0s)

```
