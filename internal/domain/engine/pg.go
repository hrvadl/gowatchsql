package engine

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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

func (e *postgreSQL) Execute(ctx context.Context, query string) error {
	if _, err := e.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("execute command: %w", err)
	}
	return nil
}

func (e *postgreSQL) GetColumns(ctx context.Context, table string) ([]Row, []Column, error) {
	const queryFmt = `
		SELECT *
		FROM information_schema.columns
    WHERE table_name   = '%s'
		`
	query := strings.TrimSpace(fmt.Sprintf(queryFmt, table))

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

func (e *postgreSQL) GetIndexes(ctx context.Context, table string) ([]Row, []Column, error) {
	const queryFmt = `
		SELECT *
		FROM pg_indexes
		WHERE tablename = '%s'
	`
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

func (e *postgreSQL) GetConstraints(ctx context.Context, table string) ([]Row, []Column, error) {
	const queryFmt = `
		SELECT conname, pg_catalog.pg_get_constraintdef(r.oid, true) as condef
		FROM pg_catalog.pg_constraint r
		WHERE r.conrelid in ('%s'::regclass)
	`
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
