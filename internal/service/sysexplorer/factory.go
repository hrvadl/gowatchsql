package sysexplorer

import (
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/hrvadl/gowatchsql/internal/platform/db"
)

func New(dsn string) (*Explorer, error) {
	switch {
	case strings.HasPrefix(dsn, "mysql://"):
		return createMySQL(cleanDBType(dsn))
	default:
		return nil, nil
	}
}

func createMySQL(dsn string) (*Explorer, error) {
	params, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := db.New("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Explorer{db, params.DBName}, nil
}

func cleanDBType(dsn string) string {
	splits := strings.Split(dsn, "://")
	if len(splits) != 2 {
		return ""
	}
	return splits[1]
}
