package sysexplorer

import (
	"context"
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
