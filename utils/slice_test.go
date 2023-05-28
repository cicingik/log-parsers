package utils

import (
	"testing"
)

func TestStringInSliceForceLower(t *testing.T) {
	data := []string{"you", "me", "her", "we", "YYYYY"}
	falseData := []string{"q", "Z", "a2bc", "xx"}

	for _, d := range data {
		exist := StringInSliceForceLower(d, data)
		if !exist {
			t.Errorf("failed find element")
		}
	}

	for _, fd := range falseData {
		exist := StringInSliceForceLower(fd, data)
		if exist {
			t.Errorf("failed find element")
		}
	}

}
