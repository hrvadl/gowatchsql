package db

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/sync/errgroup"

	"github.com/hrvadl/gowatchsql/internal/domain/errs"
)

//go:generate mockgen -destination=mocks/mock_config_repository.go -package=mocks . ConfigRepository
type ConfigRepository interface {
	AddConnection(ctx context.Context, name, dsn string) error
}

func NewPool(cfg ConfigRepository) *Pool {
	return &Pool{
		opened: make(map[string]*sqlx.DB),
		cfg:    cfg,
	}
}

type Pool struct {
	opened map[string]*sqlx.DB
	cfg    ConfigRepository
}

func (p *Pool) Get(ctx context.Context, name, driver, dsn string) (*sqlx.DB, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", errs.ErrValidation)
	}

	if driver == "" {
		return nil, fmt.Errorf("%w: driver is required", errs.ErrValidation)
	}

	if dsn == "" {
		return nil, fmt.Errorf("%w: dsn is required", errs.ErrValidation)
	}

	if conn, opened := p.opened[dsn]; opened {
		return conn, nil
	}

	conn, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}

	//	if err := conn.Ping(); err != nil {
	//		return nil, fmt.Errorf("ping: %w", err)
	//	}

	p.opened[dsn] = conn
	if err := p.cfg.AddConnection(ctx, name, dsn); err != nil {
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
