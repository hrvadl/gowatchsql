package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func New(driver, dsn string) (*sqlx.DB, error) {
	if db != nil {
		return db, nil
	}

	conn, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, err
	}

	db = conn
	return db, nil
}

func Get() *sqlx.DB {
	return db
}
