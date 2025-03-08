package engine

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type sqlite struct {
	db     *sqlx.DB
	dbPath string
}

type sqliteTable struct {
	Name string `db:"name"`
	Type string `db:"type"`
}

func (e *sqlite) Execute(ctx context.Context, query string) error {
	if _, err := e.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("execute command %q: %w", query, err)
	}
	return nil
}

func (e *sqlite) GetTables(ctx context.Context) ([]Table, error) {
	const query = `
		SELECT name, type FROM sqlite_master 
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
	`

	var tables []sqliteTable
	if err := e.db.SelectContext(ctx, &tables, query); err != nil {
		slog.Error("Got err", slog.Any("err", err))
		return nil, err
	}
	return e.toTables(tables), nil
}

func (e *sqlite) toTables(tables []sqliteTable) []Table {
	var result []Table
	slog.Debug("Converting tables", slog.Any("tables", tables))
	for _, t := range tables {
		result = append(result, Table{
			Name:   t.Name,
			Schema: "main", // SQLite doesn't use schemas in the same way as MySQL, 'main' is the default
		})
	}
	return result
}

func (e *sqlite) GetColumns(ctx context.Context, table string) ([]Row, []Column, error) {
	query := fmt.Sprintf("PRAGMA table_info('%s')", table)
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
		rowCols, err := entries.SliceScan()
		if err != nil {
			slog.Error("Got err", slog.Any("err", err))
		}
		rows = append(rows, rowCols)
	}
	if err := entries.Err(); err != nil {
		return nil, nil, err
	}
	return convertFromBinary(rows), cols, nil
}

func (e *sqlite) GetRows(ctx context.Context, table string) ([]Row, []Column, error) {
	query := fmt.Sprintf("SELECT * FROM '%s'", table)
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
		rowCols, err := entries.SliceScan()
		if err != nil {
			slog.Error("Got err", slog.Any("err", err))
		}
		rows = append(rows, rowCols)
	}
	if err := entries.Err(); err != nil {
		return nil, nil, err
	}
	return convertFromBinary(rows), cols, nil
}

func (e *sqlite) GetIndexes(ctx context.Context, table string) ([]Row, []Column, error) {
	query := fmt.Sprintf("PRAGMA index_list('%s')", table)
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
		rowCols, err := entries.SliceScan()
		if err != nil {
			slog.Error("Got err", slog.Any("err", err))
		}
		rows = append(rows, rowCols)
	}
	if err := entries.Err(); err != nil {
		return nil, nil, err
	}
	return convertFromBinary(rows), cols, nil
}

func (e *sqlite) GetConstraints(ctx context.Context, table string) ([]Row, []Column, error) {
	// SQLite doesn't have an information_schema equivalent for constraints
	// We can get foreign key constraints using PRAGMA
	query := fmt.Sprintf("PRAGMA foreign_key_list('%s')", table)
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
		rowCols, err := entries.SliceScan()
		if err != nil {
			slog.Error("Got err", slog.Any("err", err))
		}
		rows = append(rows, rowCols)
	}
	if err := entries.Err(); err != nil {
		return nil, nil, err
	}

	// We also need to check for PRIMARY KEY, UNIQUE, and CHECK constraints
	// from table_info
	pkQuery := fmt.Sprintf("PRAGMA table_info('%s')", table)
	pkEntries, err := e.db.QueryxContext(ctx, pkQuery)
	if err != nil {
		return convertFromBinary(rows), cols, nil // Return what we have so far
	}

	pkRows := make([][]any, 0)
	defer pkEntries.Close()
	for pkEntries.Next() {
		rowCols, err := pkEntries.SliceScan()
		if err != nil {
			slog.Error("Got err", slog.Any("err", err))
		}
		// Only add rows with primary key constraint (pk column > 0)
		if rowCols[5].(int64) > 0 {
			pkRows = append(pkRows, rowCols)
		}
	}

	// Combine the results
	rows = append(rows, pkRows...)

	return convertFromBinary(rows), cols, nil
}
