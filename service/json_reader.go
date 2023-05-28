package service

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cicingik/log-parsers/utils"
)

func JsonMetricReader(path string, dateRange DateRange) ([]Metric, error) {
	var metrics []Metric

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf, &metrics)
	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return []Metric{}, nil
	}

	firstMetric := metrics[0]

	t, err := firstMetric.Timestamp()
	if err != nil {
		return []Metric{}, nil
	}

	if !utils.DateInBound(t, dateRange.StartTime, dateRange.EndTime) {
		return []Metric{}, nil
	}

	return metrics, nil
}
