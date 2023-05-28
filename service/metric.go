package service

import (
	"time"

	"github.com/cicingik/log-parsers/config"
)

type Field map[string]interface{}

type Metric map[string]interface{}

type MetricsReaders func(path string, dateRange DateRange) ([]Metric, error)

type MetricInfo struct {
	FilePath string
	Metrics  []Metric
	Summary  map[string]float64
}

func (m *Metric) Timestamp() (time.Time, error) {
	metric := *m
	timestamp, err := time.Parse(time.RFC3339, metric[config.TimestampKey].(string))
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}
