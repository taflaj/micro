// access: models/database.go

package models

import (
	"database/sql"
	"net/http"

	// to isolate database details from the remainder of the application
	_ "github.com/mattn/go-sqlite3"
	"github.com/taflaj/micro/modules/database"
	"github.com/taflaj/micro/modules/messaging"
)

// DataStore implements model methods
type DataStore interface {
	Close() error
	GetAccess(int, string) (*messaging.AccessLevel, error)
	ResetAccess(http.ResponseWriter, *messaging.Message) error
	SetAccess(http.ResponseWriter, *messaging.Message) error
}

// DB is used to encapsulate sql.DB
type DB struct {
	*sql.DB
}

func initialize(db *sql.DB) error {
	// create database tables
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
	for i := 0; i < len(commands); i++ {
		if _, err = database.Execute(db, commands[i]); err != nil {
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

func (db *DB) addAddress(address int, owner string, remarks string) (int64, error) {
	command := "SELECT ID FROM Addresses WHERE Address=?"
	id, ok, err := database.GetInt(db.DB, command, address)
	if ok {
		return id, err
	}
	command = "INSERT INTO Addresses(Address, Owner, Remarks) VALUES(?, ?, ?)"
	result, err := database.Execute(db.DB, command, address, owner, remarks)
	if err != nil {
		return 0, err
	}
	id, err = result.LastInsertId()
	return id, err
}

func (db *DB) addService(service string) (int64, error) {
	command := "SELECT ID FROM Services WHERE Service=?"
	id, ok, err := database.GetInt(db.DB, command, service)
	if ok {
		return id, err
	}
	command = "INSERT INTO Services(Service) VALUES(?)"
	result, err := database.Execute(db.DB, command, service)
	if err != nil {
		return 0, err
	}
	id, err = result.LastInsertId()
	return id, err
}

func (db *DB) setAccess(address int64, service int64, canRead bool, canWrite bool) (int64, error) {
	command := "SELECT ID FROM ACL WHERE Address_ID=? AND Service_ID=?"
	id, ok, err := database.GetInt(db.DB, command, address, service)
	if ok {
		command = "UPDATE ACL\nSET CanRead=?, CanWrite=?\nWHERE ID=?"
		_, err = database.Execute(db.DB, command, canRead, canWrite, id)
		return id, err
	}
	command = "INSERT INTO ACL(Address_ID, Service_ID, CanRead, CanWrite) VALUES(?, ?, ?, ?)"
	result, err := database.Execute(db.DB, command, address, service, canRead, canWrite)
	if err != nil {
		return 0, err
	}
	id, err = result.LastInsertId()
	return id, err
}

func (db *DB) resetAccess(address int, service string) error {
	command := "DELETE FROM ACL WHERE Address_ID=(\n" +
		"  SELECT ID FROM Addresses WHERE Address=?\n" +
		") AND Service_ID=(\n" +
		"  SELECT ID FROM Services WHERE Service=?\n" +
		")"
	_, err := database.Execute(db.DB, command, address, service)
	return err
}

// GetAccess returns the access level of a given host to a given service
func (db *DB) GetAccess(address int, service string) (*messaging.AccessLevel, error) {
	level := messaging.AccessLevel{Address: address, Service: service}
	var err error
	command := "SELECT CanRead, CanWrite FROM v_ACL WHERE Address=? AND Service=?"
	row := db.DB.QueryRow(command, address, service)
	// logger.GetLogger().Printf("DEBUG row=%#v", row)
	var canRead, canWrite int
	if err = row.Scan(&canRead, &canWrite); err == nil {
		level.Defined = true
		level.CanRead = canRead != 0
		level.CanWrite = canWrite != 0
		if level.CanRead {
			level.Level = "ro"
			if level.CanWrite {
				level.Level = "rw"
			} else {
				level.Level = "no"
				if level.CanWrite {
					level.Level = "wo"
				}
			}
		}
	} else {
		level.Defined = false
	}
	return &level, err
}

// Close closes the database
func (db *DB) Close() error {
	return db.DB.Close()
}
