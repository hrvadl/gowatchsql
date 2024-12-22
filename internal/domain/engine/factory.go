package engine

import (
	"context"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	mysqlDB      = "mysql"
	postgresqlDB = "postgres"
)

type Explorer interface {
	GetTables(ctx context.Context) ([]Table, error)
	GetRows(ctx context.Context, table string) ([]Row, []Column, error)
}

type Table struct {
	Name   string `db:"TABLE_NAME"`
	Schema string `db:"TABLE_TYPE"`
}

func NewFactory(pool Pool) *Factory {
	return &Factory{
		pool: pool,
	}
}

type Pool interface {
	Get(ctx context.Context, name, driver, dsn string) (*sqlx.DB, error)
}

type Factory struct {
	pool Pool
}

func (f *Factory) Create(ctx context.Context, name, dsn string) (Explorer, error) {
	switch {
	case strings.HasPrefix(dsn, mysqlDB):
		return f.createMySQL(ctx, name, cleanDBType(dsn))
	case strings.HasPrefix(dsn, postgresqlDB):
		return f.createPostgres(ctx, name, dsn)
	default:
		return nil, errors.New("not implemented")
	}
}

func (f *Factory) createPostgres(ctx context.Context, name, dsn string) (*postgreSQL, error) {
	dsn = strings.TrimSpace(dsn)
	if !strings.Contains(dsn, "sslmode=") {
		dsn += "?sslmode=disable"
	}

	db, err := f.pool.Get(ctx, name, postgresqlDB, dsn)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(dsn, "/")
	dbName := strings.Split(parts[len(parts)-1], "?")[0]

	return &postgreSQL{db, dbName}, nil
}

func (f *Factory) createMySQL(ctx context.Context, name, dsn string) (*mySQL, error) {
	params, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	db, err := f.pool.Get(ctx, name, mysqlDB, dsn)
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
