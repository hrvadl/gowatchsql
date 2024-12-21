package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type connection struct {
	db     *sqlx.DB
	driver string
	dsn    string
}

func (c *connection) AlreadyOpened(driver, dsn string) bool {
	return c.driver == driver && c.dsn == dsn && c.db != nil
}

func (c *connection) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

var opened connection

func New(driver, dsn string) (*sqlx.DB, error) {
	if opened.AlreadyOpened(driver, dsn) {
		return opened.db, nil
	}

	conn, err := sqlx.Connect(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := opened.Close(); err != nil {
		return nil, fmt.Errorf("close old connection: %w", err)
	}

	opened = connection{conn, driver, dsn}

	return opened.db, nil
}

func Get() *sqlx.DB {
	return opened.db
}
