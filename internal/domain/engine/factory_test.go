package engine

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/hrvadl/gowatchsql/internal/domain/engine/mocks"
)

func TestNewFactory(t *testing.T) {
	t.Parallel()
	type args struct {
		pool Pool
	}
	tests := []struct {
		name string
		args args
		want *Factory
	}{
		{
			name: "Should return a new factory",
			args: args{
				pool: mocks.NewMockPool(gomock.NewController(t)),
			},
			want: &Factory{
				pool: mocks.NewMockPool(gomock.NewController(t)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := NewFactory(tt.args.pool)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFactory_Create(t *testing.T) {
	t.Parallel()
	type fields struct {
		pool func(ctrl *gomock.Controller) Pool
	}
	type args struct {
		ctx  context.Context
		name string
		dsn  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Explorer
		wantErr bool
	}{
		{
			name: "Should return an error when dsn is empty",
			fields: fields{
				pool: func(ctrl *gomock.Controller) Pool {
					return mocks.NewMockPool(ctrl)
				},
			},
			args: args{
				ctx:  t.Context(),
				name: "name",
				dsn:  "",
			},
			wantErr: true,
		},
		{
			name: "Should return an error when name is empty",
			fields: fields{
				pool: func(ctrl *gomock.Controller) Pool {
					return mocks.NewMockPool(ctrl)
				},
			},
			args: args{
				ctx:  t.Context(),
				name: "",
				dsn:  "dsn",
			},
			wantErr: true,
		},
		{
			name: "Should return an error when db type is unknown",
			fields: fields{
				pool: func(ctrl *gomock.Controller) Pool {
					return mocks.NewMockPool(ctrl)
				},
			},
			args: args{
				ctx:  t.Context(),
				name: "name",
				dsn:  "unknown",
			},
			wantErr: true,
		},
		{
			name: "Should create PosgreSQL connection",
			fields: fields{
				pool: func(ctrl *gomock.Controller) Pool {
					p := mocks.NewMockPool(ctrl)
					p.EXPECT().Get(
						gomock.Any(),
						"pg",
						postgresqlDB,
						"postgres://localhost:5432/testdb?sslmode=disable",
					).Times(1).Return(&sqlx.DB{}, nil)
					return p
				},
			},
			args: args{
				ctx:  t.Context(),
				name: "pg",
				dsn:  "postgres://localhost:5432/testdb?sslmode=disable",
			},
			want: &postgreSQL{
				schema: "testdb",
				db:     &sqlx.DB{},
			},
		},
		{
			name: "Should create MySQL connection",
			fields: fields{
				pool: func(ctrl *gomock.Controller) Pool {
					p := mocks.NewMockPool(ctrl)
					p.EXPECT().Get(
						gomock.Any(),
						"mysql",
						mysqlDB,
						"root:rootpas@(0.0.0.0:3306)/testdb",
					).Times(1).Return(&sqlx.DB{}, nil)
					return p
				},
			},
			args: args{
				ctx:  t.Context(),
				name: "mysql",
				dsn:  "mysql://root:rootpas@(0.0.0.0:3306)/testdb",
			},
			want: &mySQL{
				schema: "testdb",
				db:     &sqlx.DB{},
			},
		},
		{
			name: "Should create SQLite connection",
			fields: fields{
				pool: func(ctrl *gomock.Controller) Pool {
					p := mocks.NewMockPool(ctrl)
					p.EXPECT().Get(
						gomock.Any(),
						"lite",
						sqliteDB,
						"file.db",
					).Times(1).Return(&sqlx.DB{}, nil)
					return p
				},
			},
			args: args{
				ctx:  t.Context(),
				name: "lite",
				dsn:  "file.db",
			},
			want: &sqlite{
				dbPath: "file.db",
				db:     &sqlx.DB{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := &Factory{
				pool: tt.fields.pool(gomock.NewController(t)),
			}

			got, err := f.Create(tt.args.ctx, tt.args.name, tt.args.dsn)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
