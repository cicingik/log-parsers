package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/cicingik/log-parsers/config"
	"github.com/cicingik/log-parsers/service"
)

// command expectation
// go run main.go -d/--directory directory_name -t/--type type_file --startTime 2022-01-01T00:00:50.52Z
// --endTime 2022-01-02T00:00:00+07:00 --outputFileType[default json] json --outputFileName[default out] result

func main() {
	log.Println("start")
	start := time.Now()

	// get option and validate the value
	option, err := config.GetOption()
	if err != nil {
		flag.Usage()
		log.Printf("failed %s", err.Error())
		log.Print("kthxbye!")
		os.Exit(1)
	}

	metricService, err := service.NewMetricService(option)
	if err != nil {
		log.Printf("failed %s", err.Error())
		os.Exit(1)
	}

	// general custom filter added
	//metricService.AddReducer(service.MainMenuFilter)

	// calculate summary
	metricService.Do()

	// write report
	err = metricService.GenerateReport(option.GenerateFileName())
	if err != nil {
		log.Printf("failed %s", err.Error())
		os.Exit(1)
	}

	duration := time.Since(start)
	log.Println("done in", duration.Seconds(), "seconds")
}
