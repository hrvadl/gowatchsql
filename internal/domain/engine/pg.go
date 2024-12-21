package engine

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type postgreSQL struct {
	db     *sqlx.DB
	schema string
}

type postgreSQLTable struct {
	Name   string `db:"tablename"`
	Schema string `db:"schemaname"`
}

func (e *postgreSQL) GetTables(ctx context.Context) ([]Table, error) {
	const query = `
		SELECT tablename, schemaname FROM pg_catalog.pg_tables 
		WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'
	`

	slog.Info("Getting postgres tables", slog.Any("schema", e.schema))
	var tables []postgreSQLTable
	if err := e.db.SelectContext(ctx, &tables, query); err != nil {
		return nil, err
	}

	return e.toTables(tables), nil
}

func (e *postgreSQL) toTables(tables []postgreSQLTable) []Table {
	var t []Table
	for _, table := range tables {
		t = append(t, Table(table))
	}
	return t
}

type Column = string

type Row = []string

func (e *postgreSQL) GetRows(ctx context.Context, table string) ([]Row, []Column, error) {
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