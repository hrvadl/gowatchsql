package db

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
)

func NewPool(cfg *cfg.Config) *Pool {
	return &Pool{
		opened: make(map[string]*sqlx.DB),
		cfg:    cfg,
	}
}

type Pool struct {
	opened map[string]*sqlx.DB
	cfg    *cfg.Config
}

func (p *Pool) Get(ctx context.Context, name, driver, dsn string) (*sqlx.DB, error) {
	if conn, opened := p.opened[dsn]; opened {
		return conn, nil
	}

	conn, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}

	p.opened[dsn] = conn
	p.cfg.Connections[dsn] = name
	if err := p.cfg.Save(); err != nil {
		return nil, fmt.Errorf("save config: %w", err)
	}

	return conn, nil
}

func (p *Pool) Close() error {
	var g errgroup.Group
	for _, conn := range p.opened {
		g.Go(conn.Close)
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("close pool connections: %w", err)
	}

	return nil
}
