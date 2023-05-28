package report

import (
	"encoding/json"
	"fmt"
	"github.com/cicingik/log-parsers/config"
	"log"
	"os"
)

type Reporter func(report map[string]float64, filename string) error

func JsonReporter(summary map[string]float64, filename string) error {
	report := toReport(summary)
	b, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	// write the whole body at once
	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func YamlReporter(summary map[string]float64, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() {
		bodyCloseErr := out.Close()
		if bodyCloseErr != nil {
			log.Printf("close file error: %v", err)
		}
	}()

	// set header
	out.WriteString(fmt.Sprintf("%s\n", config.YamlHeader))
	fmt.Println(config.YamlHeader)

	for k, v := range summary {
		out.WriteString(fmt.Sprintf("- %s: %s\n", config.LevelNameKey, k))
		fmt.Println(fmt.Sprintf("- %s: %s", config.LevelNameKey, k))

		out.WriteString(fmt.Sprintf("  %s: %v\n", config.TotalValueKey, v))
		fmt.Println(fmt.Sprintf("  %s: %v", config.TotalValueKey, v))
	}

	return nil
}

func toReport(summary map[string]float64) []map[string]interface{} {
	var report []map[string]interface{}

	for k, v := range summary {
		r := map[string]interface{}{
			config.LevelNameKey:  k,
			config.TotalValueKey: v,
		}

		report = append(report, r)
	}

	return report
}
