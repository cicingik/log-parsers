package service

import (
	"time"

	"github.com/cicingik/log-parsers/config"
	"github.com/cicingik/log-parsers/utils"
)

type Reducer func(metric Metric) bool

type DateRange struct {
	// Field  is name of timestamp field
	Field string
	// StartTime is start timestamp that metric will summarize
	StartTime time.Time `json:"start_time"`
	// EndTime is end timestamp that metric will summarize
	EndTime time.Time `json:"end_time"`
}

func (i *DateRange) inBoundsMetric(metric Metric) *Metric {
	t, err := metric.Timestamp()
	if err != nil {
		return nil
	}

	if utils.InBound(t, i.StartTime, i.EndTime) {
		return &metric
	}

	return nil
}

func DateFilter(metrics []Metric, dateFilter *DateRange) []Metric {
	if dateFilter == nil {
		return metrics
	}

	chanOut := make(chan Metric)

	go func() {
		for _, metric := range metrics {
			fMetric := dateFilter.inBoundsMetric(metric)
			if fMetric != nil {
				chanOut <- *fMetric
			}
		}
		close(chanOut)
	}()

	var result []Metric
	for ch := range chanOut {
		result = append(result, ch)
	}

	return result
}

func Filter(metrics []Metric, reducer Reducer) []Metric {

	chanOut := make(chan Metric)

	go func() {
		for _, metric := range metrics {
			if reducer(metric) {
				chanOut <- metric
			}
		}
		close(chanOut)
	}()

	var result []Metric
	for ch := range chanOut {
		result = append(result, ch)
	}

	return result
}

// MainMenuFilter is sample custom filter
func MainMenuFilter(metric Metric) bool {
	if metric[config.LevelNameKey] == config.MainMenuLevel {
		return true
	}
	return false
}
