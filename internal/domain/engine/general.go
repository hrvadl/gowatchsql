package engine

import (
	"fmt"
)

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
