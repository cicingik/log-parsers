package config

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/cicingik/log-parsers/utils"
)

type (
	Option struct {
		// StartTime is start timestamp that metric will summarize
		StartTime time.Time `json:"start_time"`
		// EndTime is end timestamp that metric will summarize
		EndTime time.Time `json:"end_time"`
		// Directory is full path of directory that metric data is saved
		Directory string `json:"directory"`
		// FileType is type of data inside directory
		// only valid for json and csv
		FileType string `json:"file_type"`
		// OutputFileType is type of report summary
		// only valid for json and yaml
		OutputFileType string `json:"output_file_type"`
		// OutputFileName is name of report summary
		OutputFileName string `json:"output_file_name"`
	}
)

func GetOption() (*Option, error) {
	var (
		directory      string
		fileType       string
		startTime      string
		endTime        string
		outputFileType string
		outputFileName string
	)

	// directory argument
	flag.StringVar(&directory, "d", "", "the directory containing metric files")
	flag.StringVar(&directory, "directory", "", "the directory containing metric files")
	// file type log
	typeDesc := fmt.Sprintf("the metric file type: \n\t - %s", strings.Join(AvailableInputFileList, ",\n\t - "))
	flag.StringVar(&fileType, "t", "", typeDesc)
	flag.StringVar(&fileType, "type", "", typeDesc)
	// start time
	flag.StringVar(&startTime, "startTime", "", "the start time of the time range; RFC3339 format")
	// end time
	flag.StringVar(&endTime, "endTime", "", "the end time of the time range; RFC3339 format")
	// output type
	outTypeDesc := fmt.Sprintf("the output file type: \n\t - %s", strings.Join(AvailableOutputFileList, ",\n\t - "))
	flag.StringVar(&outputFileType, "outputFileType", "json", outTypeDesc)
	// output file name
	flag.StringVar(&outputFileName, "outputFileName", "out", "the output file type")

	// get arg from console
	flag.Parse()

	// validate
	if len(startTime) == 0 || len(endTime) == 0 {
		return nil, fmt.Errorf("missing mandatory startTime or endTime option")
	}

	st, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, err
	}

	et, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, err
	}

	//normalize OutputFileName input, remove extension if any
	name := strings.Split(strings.Trim(outputFileName, ""), ".")
	opt := Option{
		StartTime:      st,
		EndTime:        et,
		Directory:      directory,
		FileType:       fileType,
		OutputFileType: outputFileType,
		OutputFileName: name[0],
	}

	return &opt, nil

}

func (o *Option) Validate() error {
	if len(o.Directory) == 0 {
		return fmt.Errorf("missing mandatory directory option")
	}

	if len(o.FileType) == 0 {
		return fmt.Errorf("missing mandatory type option")
	}

	isDir := utils.ValidateDirectory(o.Directory)
	if !isDir {
		return fmt.Errorf("%s is not a valid directory", o.Directory)
	}

	if o.EndTime.Before(o.StartTime) {
		return fmt.Errorf("endTime must greater than startTime")
	}

	if !utils.StringInSliceForceLower(o.FileType, AvailableInputFileList) {
		return fmt.Errorf(`unsupported "%s" as input file type`, o.FileType)
	}

	if !utils.StringInSliceForceLower(o.OutputFileType, AvailableOutputFileList) {
		return fmt.Errorf(`unsupported "%s" output file type`, o.OutputFileType)
	}

	if len(o.OutputFileName) == 0 {
		o.OutputFileName = DefaultOutputName
	}

	return nil
}

func (o *Option) GenerateFileName() string {
	return fmt.Sprintf("%s.%s", o.OutputFileName, o.OutputFileType)
}
