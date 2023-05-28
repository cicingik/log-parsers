package utils

import (
	"os"
)

func ValidateDirectory(dir string) bool {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}
