package sysexplorer

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type mySQL struct {
	db     *sqlx.DB
	schema string
}

type mySQLTable struct {
	Name string `db:"TABLE_NAME"`
	Type string `db:"TABLE_TYPE"`
}

func (e *mySQL) GetTables(ctx context.Context) ([]Table, error) {
	const query = `
		SELECT TABLE_NAME, TABLE_TYPE FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_SCHEMA=? 
	`

	var tables []mySQLTable
	if err := e.db.SelectContext(ctx, &tables, query, e.schema); err != nil {
		return nil, err
	}

	return e.toTables(tables), nil
}

func (e *mySQL) toTables(tables []mySQLTable) []Table {
	var result []Table
	for _, t := range tables {
		result = append(result, Table{
			Name:   t.Name,
			Schema: e.schema,
		})
	}
	return result
}
