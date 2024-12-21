package engine

import (
	"context"
	"fmt"
	"log/slog"

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

func (e *mySQL) GetRows(ctx context.Context, table string) ([]Row, []Column, error) {
	const queryFmt = "SELECT * FROM %s"
	query := fmt.Sprintf(queryFmt, table)

	entries, err := e.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}

	rows := make([][]any, 0)
	cols, err := entries.Columns()
	if err != nil {
		return nil, nil, err
	}

	defer entries.Close()
	for entries.Next() {
		cols, err := entries.SliceScan()
		if err != nil {
			slog.Error("Got err", slog.Any("err", err))
		}

		rows = append(rows, cols)
	}

	if err := entries.Err(); err != nil {
		return nil, nil, err
	}

	return convertFromBinary(rows), cols, nil
}
