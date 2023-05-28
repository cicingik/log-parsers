package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateDirectory(t *testing.T) {
	fileName := "directory_test.go"
	fakeFile := "test.go"

	rootDir, _ := os.Getwd()
	trueCase := []string{
		rootDir,
	}
	falseCase := []string{
		filepath.Join(rootDir, fileName),
		filepath.Join(rootDir, fakeFile),
	}

	for _, el := range trueCase {
		isDir := ValidateDirectory(el)
		if isDir != true {
			t.Errorf("invalid status directory for: %v", rootDir)
		}
	}

	for _, el := range falseCase {
		isDir := ValidateDirectory(el)
		if isDir != false {
			t.Errorf("invalid status directory for: %v", rootDir)
		}
	}
}
