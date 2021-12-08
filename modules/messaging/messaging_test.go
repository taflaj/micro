// modules/messaging/messaging_test.go

package messaging

import (
	"testing"

	"github.com/taflaj/micro/modules/logger"
)

func init() {
	logger.NewLogger("Unit Test", logger.Debug)
}

func TestIPtoInt(t *testing.T) {
	var tests = []struct {
		given    string
		expected uint32
	}{
		{"70.185.37.124", 1186538876},
		{"70.188.219.55", 1186782007},
		{"127.0.0.1", 2130706433},
		{"192.168.1.1", 3232235777},
		{"192.168.1.166", 3232235942},
		{"192.168.1.194", 3232235970},
		{"192.168.1.246", 3232236022},
	}
	for _, test := range tests {
		result, err := IPtoInt(test.given)
		if err != nil {
			t.Error(err)
		} else if result != test.expected {
			t.Errorf("IPtoInt failed for %v: expected %v but got %v", test.given, test.expected, result)
		}
	}
}

func TestIPtoString(t *testing.T) {
	var tests = []struct {
		expected string
		given    int
	}{
		{"70.188.219.55", 1186782007},
		{"127.0.0.1", 2130706433},
		{"192.168.1.1", 3232235777},
		{"192.168.1.16", 3232235792},
		{"192.168.1.27", 3232235803},
	}
	for _, test := range tests {
		result := IPtoString(test.given)
		if result != test.expected {
			t.Errorf("IPtoInt failed for %v: expected %v but got %v", test.given, test.expected, result)
		}
	}
}
