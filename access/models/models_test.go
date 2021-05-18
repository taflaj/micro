// access: models/models_test.go

package models

import (
	"testing"

	"github.com/taflaj/services/modules/messaging"
)

const dsn = "/home/zezo/projects/go/src/github.com/taflaj/services/access.db"

func TestOpen(t *testing.T) {
	db, err := Open(dsn)
	if err != nil {
		t.Error(err)
	} else if db == nil {
		t.Errorf("db is %v", db)
	} else if err = db.Close(); err != nil {
		t.Error(err)
	}
}

func TestGetAccess(t *testing.T) {
	var tests = []struct {
		ip       string
		service  string
		defined  bool
		level    string
		canRead  bool
		canWrite bool
	}{
		{"127.0.0.1", "access", true, "rw", true, true},
		{"192.168.1.16", "access", true, "rw", true, true},
		{"192.168.1.27", "access", true, "rw", true, true},
		{"192.168.1.35", "digest", true, "no", false, false},
		{"192.168.1.35", "pubkey", false, "", false, false},
		{"192.168.1.35", "random", false, "", false, false},
		{"192.168.1.36", "router", false, "", false, false},
		{"192.168.1.36", "access", false, "", false, false},
		{"192.168.1.37", "digest", true, "no", false, false},
		{"192.168.1.37", "random", false, "", false, false},
	}
	db, _ := Open(dsn)
	defer db.Close()
	for _, test := range tests {
		address, _ := messaging.IPtoInt(test.ip)
		access := db.GetAccess(address, test.service)
		if test.defined == access.Defined {
			if test.defined {
				if test.level != access.Level {
					t.Errorf("%v/%v: Level should be \"%v\" but is \"%v\"", test.ip, test.service, test.level, access.Level)
				} else if test.canRead != access.CanRead {
					t.Errorf("%v/%v: CanRead should be %v but is %v", test.ip, test.service, test.canRead, access.CanRead)
				} else if test.canWrite != access.CanWrite {
					t.Errorf("%v/%v: CanWrite should be %v but is %v", test.ip, test.service, test.canWrite, access.CanWrite)
				}
			}
		} else {
			t.Errorf("%v/%v: Defined should be %v but is %v", test.ip, test.service, test.defined, access.Defined)
		}
	}
}
