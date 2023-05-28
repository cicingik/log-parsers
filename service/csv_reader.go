package service

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/cicingik/log-parsers/config"
	"github.com/cicingik/log-parsers/utils"
)

//func CSVMetricReader(path string, dateRange DateRange) ([]Metric, error) {
//	f, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func() {
//		bodyCloseErr := f.Close()
//		if bodyCloseErr != nil {
//			log.Printf("close file error: %v", err)
//		}
//	}()
//
//	csvReader := csv.NewReader(f)
//	data, err := csvReader.ReadAll()
//	if err != nil {
//		return nil, err
//	}
//	metrics := stringToMetrics(data, dateRange)
//
//	return metrics, nil
//}

//func stringToMetrics(data [][]string, dateRange DateRange) []Metric {
//	var metrics []Metric
//	for i, line := range data {
//		if i > 0 { // omit header
//
//			metric := Metric{}
//			metric[config.TimestampKey] = line[0]
//			metric[config.LevelNameKey] = line[1]
//			value, err := strconv.ParseFloat(line[2], 64)
//			if err != nil {
//				continue
//			}
//			metric[config.ValueKey] = value
//
//			if i == 1 {
//				if !utils.DateInBound(metric.Timestamp(), dateRange.StartTime, dateRange.EndTime) {
//					break
//				}
//			}
//
//			metrics = append(metrics, metric)
//		}
//	}
//	return metrics
//}

func CSVMetricReader(path string, dateRange DateRange) ([]Metric, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		bodyCloseErr := f.Close()
		if bodyCloseErr != nil {
			log.Printf("close file error: %v", err)
		}
	}()

	mettrics := []Metric{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.Split(sc.Text(), ",")

		metric := Metric{}
		metric[config.TimestampKey] = line[0]
		metric[config.LevelNameKey] = line[1]
		value, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			continue
		}
		metric[config.ValueKey] = value

		t, err := metric.Timestamp()
		if err != nil {
			continue
		}

		inBound := utils.DateInBound(t, dateRange.StartTime, dateRange.EndTime)
		if !inBound {
			break
		}

		mettrics = append(mettrics, metric)
	}
	if err := sc.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return nil, nil
	}
	return mettrics, nil
}
