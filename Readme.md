## Metric Parser

This program to summarize the data in the metric files reported daily to help the game designers to decide which level should be updated.

### Assumption

- the directory contains single type of file
- available file type input are `json`, `csv`
    
    example:

    `json:`
    ```python
    [
      {
        "timestamp": "2022-01-01T00:00:00.00Z",
        "level_name": "lobby_screen",
        "value": 73
      }
    ]
    ```
  
    `csv:`
    ```csv
    timestamp,level_name,value
    2022-01-01T00:00:00.00Z,lobby_screen,73
    ```
- 1 file contains metrics usage for 1 day; so it can assume that there is no metric with same date store in multi file
- no guarantee that there must be a file for each day
- time format `rfc3339` with `UTC` timezone
- available file type output `json`, `yaml`, `yml` (?)
- default output name file `out.{ext}`
- use standard library golang [ref.](https://pkg.go.dev/std#stdlib)

### Process

- find metric from specific date range (in `rfc3339`)
- grouping by `level_name` and sum of `value`

### Output

- display result process in the console
- write in file based output file name and type

### Available command

```shell
Usage of main.go:
  -d string
        the directory containing metric files
  -directory string
        the directory containing metric files
  -endTime string
        the end time of the time range; RFC3339 format
  -outputFileName string
        the output file type (default "out")
  -outputFileType string
        the output file type: 
                 - json,
                 - yml,
                 - yaml (default "json")
  -startTime string
        the start time of the time range; RFC3339 format
  -t string
        the metric file type: 
                 - json,
                 - csv
  -type string
        the metric file type: 
                 - json,
                 - csv
```

#### Example command
```shell
go run main.go -d ./metrics/json/ -t json --startTime 2022-02-13T00:00:50.52Z --endTime 2022-02-28T00:00:00.52Z
```

#### Output
```shell
[
    {
        "level_name": "lobby_screen",
        "total_value": 575
    },
    {
        "level_name": "level1",
        "total_value": 649
    },
    {
        "level_name": "main_menu",
        "total_value": 22
    },
    {
        "level_name": "options",
        "total_value": 60
    },
    {
        "level_name": "help",
        "total_value": 50
    }
]
```

### Use as Code
```go

import (
	...
)


type MetricService struct {
    option   *config.Option
    reader   MetricsReaders
    reporter report.Reporter
    reducers []Reducer
    Summary  map[string]float64
}

cfg := config.Option{
    StartTime:      time.Time{}, // mandatory
    EndTime:        time.Time{}, // mandatory
    Directory:      "", // mandatory
    FileType:       "", // support for json, csv
    OutputFileType: "", // support for json, yaml, yml
    OutputFileName: "",
}

metricService, _ := service.NewMetricService(option)

// set readers depend on FileType in cfg
metricService.SetReaders(func(path string) ([]Metric, error))

// writer summary depend on OutputFileType and OutputFileName
metricService.SetReporter(func(report map[string]float64, filename string) error)

// metric custom filter, can set more than one
metricService.AddReducer(func(metric Metric) bool)
metricService.AddReducer(func(metric Metric) bool)

metricService.Do()

summary := metricService.Summary

// write summary to file and show in console
_ = metricService.GenerateReport("filename.ext")

...
```

### Author

* Dany Satyanegara `<danysatyanegara@gmail.com>`

1 file contains metrics usage for 1 day, the row is sorted by timestamp
ascending (the oldest is at the top)