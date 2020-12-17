package cleaner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type PushgatewayCleaner struct {
	address string
	dryRun  bool
	ttl     time.Duration
}

func NewPushgatewayCleaner(address string, dryRun bool, ttl time.Duration) *PushgatewayCleaner {

	return &PushgatewayCleaner{address: address, dryRun: dryRun, ttl: ttl}
}

func (p *PushgatewayCleaner) collectExpired() ([]string, error) {

	urls := make([]string, 0)

	pushgatewayMetricsURL := fmt.Sprintf("%s/api/v1/metrics", p.address)
	response, err := http.Get(pushgatewayMetricsURL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	r := new(Response)
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	expiration := time.Now().Add(-p.ttl)
	logrus.Debug("Expiration date for comparison: ", expiration.UTC().String())

	for _, m := range r.Data {

		expired := m.PushTimeSeconds.Timestamp.Before(expiration)

		urlSuffix := ""

		labels := m.PushTimeSeconds.Metrics[0].Labels

		job, ok := labels["job"]
		if ok {

			urlSuffix = urlSuffix + "job" + "/" + job + "/"

			for k, v := range m.PushTimeSeconds.Metrics[0].Labels {
				if k != "" && v != "" && k != "job" {
					urlSuffix = urlSuffix + k + "/" + v + "/"
				}
			}

			test := fmt.Sprintf("%s/metrics/%s", p.address, urlSuffix)
			if expired {
				urls = append(urls, strings.TrimSuffix(test, "/"))
			}
		}

	}

	return urls, nil

}

func (p *PushgatewayCleaner) Run() error {

	logrus.Info("----------------------------")
	logrus.Info("Pushgateway Cleaner")
	logrus.Info("----------------------------")
	logrus.Info("Address:\t", p.address)
	logrus.Info("TTL:\t\t", p.ttl.String())
	logrus.Info("Dry Run:\t", p.dryRun)
	logrus.Info("----------------------------")

	logrus.Debug("Starting expired collection")
	urls, err := p.collectExpired()
	if err != nil {
		return err
	}
	logrus.Debug("Total expired count: ", len(urls))

	if p.dryRun {
		logrus.Warn("Dry run enabled, no data will be deleted!")
		logrus.Debug("The following URLs are collected:")
		for _, url := range urls {
			logrus.Debug(url)
		}
		return nil
	}

	for _, url := range urls {
		err := delete(url)
		if err != nil {
			return err
		}
	}

	logrus.Debug("Completed succesfully!")
	return nil
}

type Response struct {
	Status    string `json:"status"`
	Data      []Data `json:"data,omitempty"`
	ErrorType string `json:"errorType,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Data struct {
	PushTimeSeconds TimestampGauge `json:"push_time_seconds"`
}

type TimestampGauge struct {
	Timestamp time.Time `json:"time_stamp"`
	Metrics   []Metric  `json:"metrics"`
}

type Metric struct {
	Labels map[string]string `json:"labels"`
}

func delete(url string) error {

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("http status: %v", resp.Status)
	}

	logrus.Debug("Deletion is successful: ", url)

	return nil
}
