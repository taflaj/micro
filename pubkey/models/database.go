// pubkey: models/database.go

package models

import (
	"database/sql"

	// to isolate database details from the remainder of the application
	_ "github.com/mattn/go-sqlite3"
	"github.com/taflaj/services/modules/database"
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

func initialize(db *sql.DB) error {
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
	for _, command := range commands {
		if _, err := database.Execute(db, command); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Open connects to the database and initializes it`
func Open(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
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
