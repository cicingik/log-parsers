package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	opt := Option{
		StartTime:      time.Time{},
		EndTime:        time.Time{},
		Directory:      "",
		FileType:       "",
		OutputFileType: "",
		OutputFileName: "",
	}

	fileName := "option_test.go"

	rootDir, _ := os.Getwd()

	tMinus := time.Now().AddDate(0, 0, -1)
	tZero := time.Now()
	tOne := time.Now().AddDate(0, 0, 1)

	opt.StartTime = tZero
	opt.EndTime = tOne
	opt.Directory = rootDir

	err := opt.Validate()
	if err != nil {
		t.Errorf("validation valiled")
	}

	opt.EndTime = tMinus
	err = opt.Validate()
	if err == nil {
		t.Errorf("validation valiled")
	}

	opt.EndTime = tOne
	opt.Directory = filepath.Join(rootDir, fileName)
	err = opt.Validate()
	if err == nil {
		t.Errorf("validation valiled")
	}

}
