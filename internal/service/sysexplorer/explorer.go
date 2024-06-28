package sysexplorer

import "github.com/jmoiron/sqlx"

type Explorer struct {
	db     *sqlx.DB
	schema string
}

type Table struct {
	Name string `db:"TABLE_NAME"`
	Type string `db:"TABLE_TYPE"`
}

func (e *Explorer) GetTables() ([]Table, error) {
	const query = `
		SELECT TABLE_NAME, TABLE_TYPE FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_SCHEMA=? 
	`

	var tables []Table
	if err := e.db.Select(&tables, query, e.schema); err != nil {
		return nil, err
	}

	return tables, nil
}
