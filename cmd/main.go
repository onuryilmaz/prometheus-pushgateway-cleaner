package main

import (
	"flag"
	"time"

	"github.com/onuryilmaz/prometheus-pushgateway-cleaner/pkg/cleaner"
	"github.com/sirupsen/logrus"
)

var (
	ttl     = flag.Duration("ttl", 24*time.Hour, "TTL for clearing the expired metrics, the default is 24 hours.")
	dryrun  = flag.Bool("dry-run", false, "Dry run and do not delete the metrics, the default is false")
	address = flag.String("address", "", "Address of Prometheus Pushgateway")
	debug   = flag.Bool("debug", false, "Set debug log level, the default is false.")
)

func main() {

	flag.Parse()

	if *address == "" {
		logrus.Fatal("address should be provided")
	}

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	c := cleaner.NewPushgatewayCleaner(*address, *dryrun, *ttl)
	c.Run()
}
