package db

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// @TODO: save contexts.
func NewPool() *Pool {
	return &Pool{}
}

type Pool struct {
	opened connection
}

func (p *Pool) GetOrCreate(ctx context.Context, driver, dsn string) (*sqlx.DB, error) {
	if p.opened.AlreadyOpened(driver, dsn) {
		return p.opened.db, nil
	}

	conn, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := p.opened.Close(); err != nil {
		return nil, fmt.Errorf("close old connection: %w", err)
	}

	p.opened = connection{conn, driver, dsn}

	return p.opened.db, nil
}

func (p *Pool) Get() *sqlx.DB {
	return p.opened.db
}

func (p *Pool) Close() error {
	return p.opened.Close()
}
