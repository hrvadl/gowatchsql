package db

import "github.com/jmoiron/sqlx"

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
