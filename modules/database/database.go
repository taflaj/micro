// modules/database/database.go

package database

import (
	"database/sql"
)

// Execute executes a SQL statement
func Execute(db *sql.DB, cmd string, params ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(params...)
}

// GetInt performs a query and returns a single integer
func GetInt(db *sql.DB, query string, params ...interface{}) (int64, bool, error) {
	row := db.QueryRow(query, params...)
	var result int64
	if err := row.Scan(&result); err != nil {
		return 0, false, err
	}
	return result, true, nil
}
