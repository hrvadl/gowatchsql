package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/hrvadl/gowatchsql/internal/domain/errs"
)

const (
	mysqlDB      = "mysql"
	postgresqlDB = "postgres"
	fileDBSuffix = ".db"
	sqliteDB     = "sqlite3"
)

type Explorer interface {
	GetTables(ctx context.Context) ([]Table, error)
	GetRows(ctx context.Context, table string) ([]Row, []Column, error)
	GetColumns(ctx context.Context, table string) ([]Row, []Column, error)
	GetIndexes(ctx context.Context, table string) ([]Row, []Column, error)
	GetConstraints(ctx context.Context, table string) ([]Row, []Column, error)
	Execute(ctx context.Context, query string) error
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

//go:generate mockgen -destination=mocks/mock_pool.go -package=mocks . Pool
type Pool interface {
	Get(ctx context.Context, name, driver, dsn string) (*sqlx.DB, error)
}

type Factory struct {
	pool Pool
}

func (f *Factory) Create(ctx context.Context, name, dsn string) (Explorer, error) {
	if dsn == "" {
		return nil, fmt.Errorf("%w: dsn is required", errs.ErrValidation)
	}

	if name == "" {
		return nil, fmt.Errorf("%w: name is required", errs.ErrValidation)
	}

	switch {
	case strings.HasPrefix(dsn, postgresqlDB):
		return f.createPostgres(ctx, name, dsn)
	case strings.Contains(dsn, fileDBSuffix):
		return f.createSQLite(ctx, name, dsn)
	case strings.HasPrefix(dsn, mysqlDB) || !strings.HasPrefix(dsn, "://"):
		return f.createMySQL(ctx, name, cleanDBType(dsn))
	default:
		return nil, fmt.Errorf("%w: unsupported database type", errs.ErrValidation)
	}
}

func (f *Factory) createSQLite(ctx context.Context, name, file string) (*sqlite, error) {
	db, err := f.pool.Get(ctx, name, sqliteDB, file)
	if err != nil {
		return nil, fmt.Errorf("connect to sqlite: %w", err)
	}

	return &sqlite{db, file}, nil
}

func (f *Factory) createPostgres(ctx context.Context, name, dsn string) (*postgreSQL, error) {
	dsn = strings.TrimSpace(dsn)
	if !strings.Contains(dsn, "sslmode=") {
		dsn += "?sslmode=disable"
	}

	db, err := f.pool.Get(ctx, name, postgresqlDB, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	parts := strings.Split(dsn, "/")
	dbName := strings.Split(parts[len(parts)-1], "?")[0]

	return &postgreSQL{db, dbName}, nil
}

func (f *Factory) createMySQL(ctx context.Context, name, dsn string) (*mySQL, error) {
	params, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("validate mysql url: %w", err)
	}

	db, err := f.pool.Get(ctx, name, mysqlDB, dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql: %w", err)
	}

	return &mySQL{db, params.DBName}, nil
}

func cleanDBType(dsn string) string {
	splits := strings.Split(dsn, "://")
	if len(splits) != 2 {
		return dsn
	}
	return splits[1]
}
