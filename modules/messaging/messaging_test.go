// modules/messaging/messaging_test.go

package messaging

import (
	"testing"

	"github.com/taflaj/services/modules/logger"
)

func init() {
	logger.NewLogger("Unit Test", logger.Debug)
}

func TestIPtoInt(t *testing.T) {
	var tests = []struct {
		given    string
		expected int
	}{
		{"70.188.219.55", 1186782007},
		{"127.0.0.1", 2130706433},
		{"192.168.1.1", 3232235777},
		{"192.168.1.16", 3232235792},
		{"192.168.1.27", 3232235803},
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

func TestGetData(t *testing.T) {
	var tests = []struct {
		message  Message
		expected Map
	}{
		{Message{IP: 0}, Map{}},
		{Message{IP: 1, Data: []string{""}}, Map{"": ""}},
		{Message{IP: 2, Data: []string{"ip=127.0.0.1"}}, Map{"ip": "127.0.0.1"}},
		{Message{IP: 3, Data: []string{"ip=192.168.1.11", "one=1", "two=2"}}, Map{"ip": "192.168.1.11", "one": "1", "two": "2"}},
		{Message{IP: 4, Data: []string{"ip=192.145.119.180", "three=3", "four=4", "etc."}}, Map{"ip": "192.145.119.180", "three": "3", "four": "4", "etc.": ""}},
	}
	for _, test := range tests {
		result := test.message.GetData()
		if test.expected == nil {
			if result != nil {
				t.Errorf("GetData failed for %v: expected nil but got %v", test.message.IP, result)
			}
		} else if result == nil {
			t.Errorf("GetData failed for %v: expected a Map but got nil", test.message.IP)
		} else {
			le := len(test.expected)
			lr := len(*result)
			if le != lr {
				t.Errorf("GetData failed for %v: expected a Map of length %v but got %v", test.message.IP, le, lr)
			} else {
				for k, v := range test.expected {
					r, ok := (*result)[k]
					if !ok {
						t.Errorf("GetData failed for %v: unable for find element \"%v\"", test.message.IP, k)
					} else if v != r {
						t.Errorf("GetData failed for %v: expected element \"%v\" to be \"%v\" but got \"%v\"", test.message.IP, k, v, r)
					}
				}
			}
		}
	}
}
