package lib

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

// Connects to db -> defer when instantiating
func ConnectDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) CreateTable(table_name string, structure string) error {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", table_name, structure)
	_, err := db.db.Exec(query)
	return err

}

// need to defer rows.Close() right after calling the function in main?
func (db *DB) SelectFromTable(table string, columns string, _query ...string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", columns, table)
	if len(_query) == 1 {
		query = fmt.Sprintf("%s WHERE %s", query, _query)
	}
	rows, err := db.db.Query(query + ";")
	if err != nil {
		return nil, err
	}
	// defer rows.Close()
	return rows, nil
}

func (db *DB) InsertRecordIntoTable(table string, data map[string]interface{}) error {
	values := "("
	for _, value := range data {
		switch reflect.ValueOf(value).Kind() {

		case reflect.String:
			values = fmt.Sprintf("%s'%s', ", values, value)

		default:
			values = fmt.Sprintf("%s%d, ", values, value)
		}
	}
	values = values[:len(values)-2] + ")"

	query := fmt.Sprintf("INSERT OR IGNORE INTO %s VALUES %s;", table, values)
	_, err := db.db.Exec(query)
	return err

}

func (db *DB) UpdateRecord(table string, data map[string]interface{}, _query string) error {
	values := "("
	for key, value := range data {
		switch reflect.ValueOf(value).Kind() {

		case reflect.String:
			values = fmt.Sprintf("%s%s = '%s', ", values, key, value)

		default:
			values = fmt.Sprintf("%s%s = %s, ", values, key, value)
		}
	}
	values = values[2:] + ")"
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", table, values, _query)
	_, err := db.db.Exec(query)
	return err

}

func (db *DB) DeleteFromTable(table string, key string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s;", table, key)
	_, err := db.db.Exec(query)
	return err
}
