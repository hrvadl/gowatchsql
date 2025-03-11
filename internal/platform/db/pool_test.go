package db

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/gowatchsql/internal/platform/db/mocks"
	"github.com/hrvadl/gowatchsql/pkg/xtest"
)

func TestNewPool(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	t.Parallel()
	type args struct {
		cfg ConfigRepository
	}
	tests := []struct {
		name string
		args args
		want *Pool
	}{
		{
			name: "Should return a new pool",
			args: args{
				cfg: mocks.NewMockConfigRepository(gomock.NewController(t)),
			},
			want: &Pool{
				cfg:    mocks.NewMockConfigRepository(gomock.NewController(t)),
				opened: make(map[string]*sqlx.DB),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewPool(tt.args.cfg)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestPool_Close(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	t.Parallel()
	type fields struct {
		opened map[string]*sqlx.DB
		cfg    ConfigRepository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Should close all connections",
			fields: fields{
				opened: make(map[string]*sqlx.DB),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &Pool{
				opened: tt.fields.opened,
				cfg:    tt.fields.cfg,
			}

			require.NoError(t, p.Close())
		})
	}
}

func TestPool_Get(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	t.Parallel()
	type fields struct {
		opened map[string]*sqlx.DB
		cfg    ConfigRepository
	}
	type args struct {
		ctx    context.Context
		name   string
		driver string
		dsn    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sqlx.DB
		wantErr bool
	}{
		{
			name: "Should return cached connection",
			fields: fields{
				opened: map[string]*sqlx.DB{
					"dsn": {},
				},
			},
			args: args{
				ctx:    t.Context(),
				name:   "name",
				dsn:    "dsn",
				driver: "driver",
			},
			want: &sqlx.DB{},
		},
		{
			name:   "Should return an error when name is empty",
			fields: fields{},
			args: args{
				ctx:    t.Context(),
				name:   "",
				dsn:    "dsn",
				driver: "driver",
			},
			wantErr: true,
		},
		{
			name:   "Should return an error when dsn is empty",
			fields: fields{},
			args: args{
				ctx:    t.Context(),
				name:   "name",
				dsn:    "",
				driver: "driver",
			},
			wantErr: true,
		},
		{
			name:   "Should return an error when driver is empty",
			fields: fields{},
			args: args{
				ctx:    t.Context(),
				name:   "name",
				dsn:    "dsn",
				driver: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := &Pool{
				opened: tt.fields.opened,
			}

			got, err := p.Get(tt.args.ctx, tt.args.name, tt.args.driver, tt.args.dsn)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
