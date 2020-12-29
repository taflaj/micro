// pubkey: models/database.go

package models

import (
	"database/sql"

	// to isolate database details from the remainder of the application
	_ "github.com/mattn/go-sqlite3"
)

// DataStore implements model methods
type DataStore interface {
	Close() error
	GetPublicKey(string) (string, error)
}

// DB is used to encapsulate sql.DB
type DB struct {
	*sql.DB
}

func execute(db *sql.DB, cmd string, params ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(params...)
}

func initialize(db *sql.DB) error {
	// create database tables
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	commands := []string{
		"CREATE TABLE IF NOT EXISTS PubKeys(\n" +
			"  ID INTEGER PRIMARY KEY AUTOINCREMENT,\n" +
			"  Email TEXT NOT NULL,\n" +
			"  PubKey TEXT)",
		"CREATE UNIQUE INDEX IF NOT EXISTS PubKeys_ID_IDX ON PubKeys(ID, Email)",
	}
	for i := 0; i < len(commands); i++ {
		if _, err = execute(db, commands[i]); err != nil {
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

// Close closes the database
func (db *DB) Close() error {
	return db.DB.Close()
}

// GetPublicKey returns a public key based on the email address
func (db *DB) GetPublicKey(email string) (string, error) {
	query := "SELECT PubKey FROM PubKeys WHERE Email=?"
	row := db.QueryRow(query, email)
	var key string
	if err := row.Scan(&key); err != nil {
		return "", err
	}
	return key, nil
}
