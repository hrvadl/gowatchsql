package tableexplorer

import "github.com/jmoiron/sqlx"

func New() *Explorer {
	return &Explorer{}
}

type Explorer struct {
	db    *sqlx.DB
	table string
}

func (e *Explorer) GetAll() (map[string]any, error) {
	const query = "SELECT * FROM ?"
	dst := make(map[string]any)
	if err := e.db.Select(&dst, query, e.table); err != nil {
		return nil, err
	}
	return dst, nil
}
