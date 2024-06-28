package tableexplorer

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *Explorer {
	return &Explorer{
		db: db,
	}
}

type Explorer struct {
	db *sqlx.DB
}

type Column = string

type Row = []string

func (e *Explorer) GetAll(table string) ([]Row, []Column, error) {
	const queryFmt = "SELECT * FROM %s"
	query := fmt.Sprintf(queryFmt, table)

	entries, err := e.db.Queryx(query)
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

func convertFromBinary(entries [][]any) []Row {
	rows := make([]Row, 0)

	for _, m := range entries {
		row := make(Row, 0)
		for _, v := range m {
			switch v := v.(type) {
			case []byte:
				row = append(row, string(v))
			default:
				row = append(row, fmt.Sprintf("%v", v))
			}
		}
		rows = append(rows, row)
	}

	return rows
}
