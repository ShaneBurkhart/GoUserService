package main

import (
	"database/sql"
)

// Up is executed when this migration is applied
func Up_20150114025706(txn *sql.Tx) {
	txn.Exec(`
		CREATE TABLE users (
		id SERIAL NOT NULL PRIMARY KEY,
		email varchar(255) UNIQUE,
		password_digest varchar(60)
		)
	`)
}

// Down is executed when this migration is rolled back
func Down_20150114025706(txn *sql.Tx) {
	txn.Exec(`
		DROP TABLE users
	`)
}
