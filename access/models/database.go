// access: models/database.go

package models

import (
	"database/sql"

	// to isolate database details from the remainder of the application
	_ "github.com/mattn/go-sqlite3"
	"github.com/taflaj/services/modules/database"
	"github.com/taflaj/services/modules/logger"
	"github.com/taflaj/services/modules/messaging"
)

// DataStore implements model methods
type DataStore interface {
	Close() error
	GetAccess(int, string) *messaging.AccessLevel
	ResetAccess(int, string) error
	SetAccess(*messaging.AccessLevel, string, string) error
}

// DB is used to encapsulate sql.DB
type DB struct {
	*sql.DB
}

func initialize(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	commands := []string{
		"CREATE TABLE IF NOT EXISTS Addresses(\n" +
			"  ID INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
			"  Address INTEGER NOT NULL,\n" +
			"  Owner TEXT,\n" +
			"  Remarks TEXT)",
		"CREATE UNIQUE INDEX IF NOT EXISTS Addresses_ID_IDX ON Addresses(ID,Address)",
		"CREATE VIEW IF NOT EXISTS v_Addresses AS\n" +
			"  SELECT ID, Address, Address>>24 AS A, Address>>16&255 AS B, Address>>8&255 AS C, Address&255 AS D, Owner, Remarks\n" +
			"  FROM Addresses",
		"CREATE TABLE IF NOT EXISTS Services(\n" +
			"  ID INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
			"  Service TEXT NOT NULL,\n" +
			"  Remarks TEXT)",
		"CREATE UNIQUE INDEX IF NOT EXISTS Services_ID_IDX ON Services(ID,Service)",
		"CREATE TABLE IF NOT EXISTS ACL(\n" +
			"  ID INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
			"  Address_ID INTEGER NOT NULL,\n" +
			"  Service_ID INTEGER NOT NULL,\n" +
			"  CanRead INTEGER NOT NULL,\n" +
			"  CanWrite INTEGER NOT NULL,\n" +
			"  Remarks TEXT)",
		"CREATE UNIQUE INDEX IF NOT EXISTS ACL_ID_IDX ON ACL(ID,Address_ID,Service_ID)",
		"CREATE VIEW IF NOT EXISTS v_ACL AS\n" +
			"  SELECT A.ID, B.Address, B.A, B.B, B.C, B.D, B.Owner, (SELECT Service FROM Services C WHERE A.Service_ID = C.ID) Service, A.CanRead, A.CanWrite\n" +
			"  FROM ACL A, v_Addresses B\n" +
			"  WHERE A.Address_ID = B.ID",
	}
	for _, command := range commands {
		if _, err = database.Execute(db, command); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Open connects to the database and initializes it
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	if err = initialize(db); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// ResetAccess restores access for a given host to a specific service to the default
func (db *DB) ResetAccess(address int, service string) error {
	command := "DELETE FROM ACL WHERE Address_ID=(\n" +
		"  SELECT ID FROM Addresses WHERE Address=?\n" +
		") AND Service_ID=(\n" +
		"  SELECT ID FROM Services WHERE Service=?\n" +
		")"
	_, err := database.Execute(db.DB, command, address, service)
	return err
}

// SetAccess defines access for a given host to a specific service
func (db *DB) SetAccess(level *messaging.AccessLevel, owner string, remarks string) error {
	if _, ok, _ := database.GetInt(db.DB, "SELECT ID FROM Addresses WHERE Address=?", level.IP); !ok {
		if _, err := database.Execute(db.DB, "INSERT INTO Addresses(Address, Owner, Remarks) VALUES(?, ?, ?)", level.IP, owner, remarks); err != nil {
			logger.GetLogger().Log(logger.Error, err)
			return err
		}
	}
	if _, ok, _ := database.GetInt(db.DB, "SELECT ID FROM Services WHERE Service=?", level.Service); !ok {
		if _, err := database.Execute(db.DB, "INSERT INTO Services(Service, Remarks) VALUES(?, ?)", level.Service, remarks); err != nil {
			logger.GetLogger().Log(logger.Error, err)
			return err
		}
	}
	id, ok, _ := database.GetInt(db.DB, "SELECT ID FROM v_ACL WHERE Address=? AND Service=?", level.IP, level.Service)
	if ok {
		if _, err := database.Execute(db.DB, "UPDATE ACL\nSET CanRead=?, CanWrite=?, Remarks=?\nWHERE ID=?", level.CanRead, level.CanWrite, remarks, id); err != nil {
			logger.GetLogger().Log(logger.Error, err)
			return err
		}
	} else {
		if _, err := database.Execute(db.DB, "INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite, Remarks) VALUES((SELECT ID FROM Addresses WHERE Address=?), (SELECT ID FROM Services WHERE Service=?), ?, ?, ?)", level.IP, level.Service, level.CanRead, level.CanWrite, remarks); err != nil {
			logger.GetLogger().Log(logger.Error, err)
			return err
		}
	}
	return nil
}

// GetAccess returns the access level of a given host to a given service
func (db *DB) GetAccess(address int, service string) *messaging.AccessLevel {
	level := messaging.AccessLevel{IP: address, Service: service}
	command := "SELECT CanRead, CanWrite FROM v_ACL WHERE Address=? AND Service=?"
	row := db.DB.QueryRow(command, address, service)
	var canRead, canWrite int
	if err := row.Scan(&canRead, &canWrite); err == nil {
		level.Defined = true
		level.CanRead = canRead != 0
		level.CanWrite = canWrite != 0
		if level.CanRead {
			level.Level = "ro"
			if level.CanWrite {
				level.Level = "rw"
			}
		} else {
			level.Level = "no"
			if level.CanWrite {
				level.Level = "wo"
			}
		}
	}
	return &level
}

// Close closes the database
func (db *DB) Close() error {
	return db.DB.Close()
}
