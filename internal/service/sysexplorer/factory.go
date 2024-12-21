package sysexplorer

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/hrvadl/gowatchsql/internal/platform/db"
)

const (
	mysqlDB      = "mysql"
	postgresqlDB = "postgres"
)

type Explorer interface {
	GetTables(context.Context) ([]Table, error)
}

type Table struct {
	Name   string `db:"TABLE_NAME"`
	Schema string `db:"TABLE_TYPE"`
}

func New(dsn string) (Explorer, error) {
	switch {
	case strings.HasPrefix(dsn, mysqlDB):
		return newMySQL(cleanDBType(dsn))
	case strings.HasPrefix(dsn, postgresqlDB):
		return newPostgres(dsn)
	default:
		return nil, errors.New("not implemented")
	}
}

func newPostgres(dsn string) (*postgreSQL, error) {
	dsn = strings.TrimSpace(dsn)
	if !strings.Contains(dsn, "sslmode=") {
		dsn += "?sslmode=disable"
	}

	db, err := db.New(postgresqlDB, dsn)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(dsn, "/")
	dbName := strings.Split(parts[len(parts)-1], "?")[0]

	return &postgreSQL{db, dbName}, nil
}

func newMySQL(dsn string) (*mySQL, error) {
	params, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := db.New(mysqlDB, dsn)
	if err != nil {
		return nil, err
	}

	return &mySQL{db, params.DBName}, nil
}

func cleanDBType(dsn string) string {
	splits := strings.Split(dsn, "://")
	if len(splits) != 2 {
		return ""
	}
	return splits[1]
}
