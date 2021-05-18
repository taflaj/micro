// modules/logger/logger_test.go

package logger

import "testing"

func TestLevel(t *testing.T) {
	var tests = []struct {
		level Level
		name  string
	}{
		{-1, "NONE"},
		{None, "NONE"},
		{1, "CUSTOM01"},
		{Debug, "DEBUG"},
		{Info, "INFO"},
		{Warning, "WARNING"},
		{Error, "ERROR"},
		{Critical, "CRITICAL"},
		{90, "CUSTOM90"},
		{100, "CUSTOM99"},
		{1000, "CUSTOM99"},
	}
	for _, test := range tests {
		NewLogger("test", test.level)
		level := GetLogger().GetLevel()
		if level < 0 && test.level != 0 {
			t.Error("Minimum level should be adjusted to 0")
		} else if level >= 100 && test.level != 99 {
			t.Error("Maximum level should be capped at 99")
		} else if level != test.level {
			if test.level < 0 && level == 0 {

			} else if test.level >= 100 && level == 99 {

			} else {
				t.Errorf("Level set to %v but got %v", test.level, level)
			}
		} else {
			name := level.ToString()
			if name != test.name {
				t.Errorf("Level %v name should be %s but got %s", level, test.name, name)
			}
		}
	}
}
