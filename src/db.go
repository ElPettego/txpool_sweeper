package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func newDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

// need to defer rows.Close() right after calling the function in main
func (db *DB) selectFromTable(table string, columns string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s;", columns, table)
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	// defer rows.Close()
	return rows, nil
}

func (db *DB) insertRecordIntoTable(table string, data interface{}) error {
	return nil

}

func (db *DB) deleteFromTable(table string, key string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s;", table, key)
	_, err := db.db.Exec(query)
	return err
}
