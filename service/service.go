package service

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/cicingik/log-parsers/config"
	"github.com/cicingik/log-parsers/report"
)

type MetricService struct {
	option   *config.Option
	reader   MetricsReaders
	reporter report.Reporter
	reducers []Reducer
	Summary  map[string]float64
}

func NewMetricService(option *config.Option) (*MetricService, error) {
	err := option.Validate()
	if err != nil {
		return nil, err
	}

	svc := MetricService{
		option: option,
	}

	// set engine service
	switch inputType := option.FileType; inputType {
	case config.JsonFile:
		svc.SetReaders(JsonMetricReader)
	case config.CsvFile:
		svc.SetReaders(CSVMetricReader)
	default:
	}

	//set report generator
	switch outputType := option.OutputFileType; outputType {
	case config.JsonFile:
		svc.SetReporter(report.JsonReporter)
	case config.YmlFile:
		svc.SetReporter(report.YamlReporter)
	case config.YamlFile:
		svc.SetReporter(report.YamlReporter)
	default:
		svc.SetReporter(report.JsonReporter)
	}

	return &svc, nil
}

func (s *MetricService) SetReaders(reader MetricsReaders) {
	s.reader = reader
}

func (s *MetricService) SetReporter(reporter report.Reporter) {
	s.reporter = reporter
}

func (s *MetricService) AddReducer(reducer Reducer) {
	s.reducers = append(s.reducers, reducer)
}

func (s *MetricService) Do() {
	dateRange := DateRange{
		Field:     config.TimestampKey,
		StartTime: s.option.StartTime,
		EndTime:   s.option.EndTime,
	}
	//read all file
	chanMetrics := s.read(dateRange)

	// TODO: make filter date can be general and can be set as Reducer
	// reduce by filter by date, will not run if DateRange set by `nil`
	chanMetrics = s.filterDate(chanMetrics, &dateRange)

	// custom filter
	chanMetrics = s.filter(chanMetrics)

	//get summary
	summaries := s.calculateSummary(chanMetrics)

	s.Summary = summaries
}

func (s *MetricService) GenerateReport(filename string) error {
	err := s.reporter(s.Summary, filename)
	if err != nil {
		return err
	}

	return nil
}

func (s *MetricService) read(dateRange DateRange) <-chan MetricInfo {
	chanOut := make(chan MetricInfo)

	go func() {
		err := filepath.Walk(s.option.Directory,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				// TODO: check if file type is same with config option

				metrics, err := s.reader(path, dateRange)
				if err != nil {
					log.Printf("failed read file %s, %s \n", path, err.Error())
				}

				chanOut <- MetricInfo{
					FilePath: path,
					Metrics:  metrics,
				}

				return nil
			})
		if err != nil {
			log.Println("failed read file", err.Error())
		}

		close(chanOut)
	}()

	return chanOut
}

func (s *MetricService) filterDate(chanIn <-chan MetricInfo, filterDate *DateRange) <-chan MetricInfo {
	chanOut := make(chan MetricInfo)

	go func() {
		for metricInfo := range chanIn {

			metrics := DateFilter(metricInfo.Metrics, filterDate)

			metricInfo.Metrics = metrics

			chanOut <- metricInfo
		}
		close(chanOut)
	}()

	return chanOut
}

func (s *MetricService) filter(chanIn <-chan MetricInfo) <-chan MetricInfo {
	chanOut := make(chan MetricInfo)

	go func() {
		for metricInfo := range chanIn {

			metrics := metricInfo.Metrics
			for _, reducer := range s.reducers {
				metrics = Filter(metrics, reducer)
			}
			metricInfo.Metrics = metrics

			chanOut <- metricInfo
		}
		close(chanOut)
	}()

	return chanOut
}

func (s *MetricService) calculateSummary(chanIn <-chan MetricInfo) map[string]float64 {
	chanOut := make(chan MetricInfo)

	go func() {
		for metricInfo := range chanIn {
			result := map[string]float64{}
			for _, metric := range metricInfo.Metrics {
				levelName, _ := metric[config.LevelNameKey]
				intVal := metric[config.ValueKey].(float64)
				stringLevelName := levelName.(string)

				_, ok := result[stringLevelName]
				if ok {
					result[stringLevelName] += intVal
				} else {
					result[stringLevelName] = intVal
				}
			}
			metricInfo.Summary = result
			chanOut <- metricInfo
		}
		close(chanOut)
	}()

	finalSummary := map[string]float64{}
	for metricInfo := range chanOut {
		for k, v := range metricInfo.Summary {
			_, ok := finalSummary[k]
			if ok {
				finalSummary[k] += v
			} else {
				finalSummary[k] += v
			}
		}
	}

	return finalSummary
}

func mergeChanMetricInfo(chanInMany ...<-chan MetricInfo) <-chan MetricInfo {
	wg := new(sync.WaitGroup)
	chanOut := make(chan MetricInfo)

	wg.Add(len(chanInMany))
	for _, eachChan := range chanInMany {
		go func(eachChan <-chan MetricInfo) {
			for eachChanData := range eachChan {
				chanOut <- eachChanData
			}
			wg.Done()
		}(eachChan)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}
