package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type MyDB struct {
	db *sql.DB
}

// ConnectDB connects to database and return a pointer to it.
func ConnectDB(dbName string, user string, pw string) (*sql.DB, error) {
	dbName = "dbname=" + dbName + " user=" + user + " password=" + pw + " " + " sslmode=disable"
	db, err := sql.Open("postgres", dbName)
	return db, err
}

// initialize with some data
func (mydb *appContext) InitDB() {
	db := mydb.db
	if _, err := setupTables(db); err != nil {
		panic(err)
	}
}

// setupTables creates the DB schema for storing metadata
func setupTables(db *sql.DB) (sql.Result, error) {
	//FIXME remove drop tables

	result, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS bleot (uuid uuid, bleotid smallint)",
	)
	return result, err
}
