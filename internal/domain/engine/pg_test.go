package engine

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/hrvadl/gowatchsql/pkg/xtest"
)

var postgresTestDSN string

func Test_postgreSQL_Execute(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	type args struct {
		ctx   context.Context
		query string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should execute query",
			args: args{
				ctx:   t.Context(),
				query: "DELETE FROM users",
			},
			wantErr: false,
		},
		{
			name: "Should return err if query is invalid",
			args: args{
				ctx:   t.Context(),
				query: "OIGHDSOIHGOIDSHGOIHS",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedPostgreSQL(t, postgresqlDB, postgresTestDSN)
			t.Cleanup(cleanup)

			e := &postgreSQL{
				db:     db,
				schema: dbName,
			}

			err := e.Execute(tt.args.ctx, tt.args.query)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func Test_postgreSQL_GetTables(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	type args struct {
		ctx    context.Context
		schema string
	}
	tests := []struct {
		name    string
		args    args
		want    []Table
		wantErr bool
	}{
		{
			name: "Should return tables",
			args: args{
				ctx:    t.Context(),
				schema: dbName,
			},
			wantErr: false,
			want: []Table{
				{Name: tableName, Schema: dbName},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedPostgreSQL(t, postgresqlDB, postgresTestDSN)
			t.Cleanup(cleanup)

			e := &postgreSQL{
				db:     db,
				schema: tt.args.schema,
			}

			got, err := e.GetTables(tt.args.ctx)
			if tt.wantErr {
				require.Len(t, got, 0)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, 1)
			table := got[0]
			require.Equal(t, tableName, table.Name)
		})
	}
}

func Test_postgreSQL_GetColumns(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	type args struct {
		ctx   context.Context
		table string
	}
	tests := []struct {
		name        string
		args        args
		wantRows    []Row
		wantColumns []Column
		wantErr     bool
	}{
		{
			name: "Should return columns",
			args: args{
				ctx:   t.Context(),
				table: tableName,
			},
		},
		{
			name: "Should return error if table is not found",
			args: args{
				ctx:   t.Context(),
				table: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedPostgreSQL(t, postgresqlDB, postgresTestDSN)
			t.Cleanup(cleanup)

			e := &postgreSQL{
				db:     db,
				schema: dbName,
			}

			gotRows, _, err := e.GetColumns(tt.args.ctx, tt.args.table)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			columnsFreq := map[string]int{
				"email":      0,
				"name":       0,
				"id":         0,
				"created_at": 0,
			}

			for col := range columnsFreq {
				var contains bool
				for _, row := range gotRows {
					if slices.Contains(row, col) {
						contains = true
						break
					}
				}
				require.Truef(t, contains, "Column %s not found", col)
			}
		})
	}
}

func Test_postgreSQL_GetRows(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	type args struct {
		ctx   context.Context
		table string
	}
	tests := []struct {
		name        string
		args        args
		wantRows    []Row
		wantColumns []Column
		wantErr     bool
	}{
		{
			name: "Should get rows",
			args: args{
				ctx:   t.Context(),
				table: tableName,
			},
			wantRows: []Row{
				{"1", "John Doe", "john@example.com", "2023-01-01 10:00:00 +0000 +0000"},
				{"2", "Jane Smith", "jane@example.com", "2023-01-01 10:00:00 +0000 +0000"},
				{"3", "Bob Wilson", "bob@example.com", "2023-01-01 10:00:00 +0000 +0000"},
			},
			wantColumns: []Column{"id", "name", "email", "created_at"},
		},
		{
			name: "Should not get rows if table is not found",
			args: args{
				ctx:   t.Context(),
				table: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedPostgreSQL(t, postgresqlDB, postgresTestDSN)
			t.Cleanup(cleanup)

			e := &postgreSQL{
				db:     db,
				schema: dbName,
			}

			gotRows, gotColumns, err := e.GetRows(tt.args.ctx, tt.args.table)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantRows, gotRows)
			require.Equal(t, tt.wantColumns, gotColumns)
		})
	}
}

func Test_postgreSQL_GetIndexes(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	type args struct {
		ctx   context.Context
		table string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should get indexes",
			args: args{
				ctx:   t.Context(),
				table: tableName,
			},
		},
		{
			name: "Should return error if table is not found",
			args: args{
				ctx:   t.Context(),
				table: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedPostgreSQL(t, postgresqlDB, postgresTestDSN)
			t.Cleanup(cleanup)

			e := &postgreSQL{
				db:     db,
				schema: dbName,
			}

			gotRows, _, err := e.GetIndexes(tt.args.ctx, tt.args.table)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.NotEmpty(t, gotRows, "Indexes not found")
		})
	}
}

func Test_postgreSQL_GetConstraints(t *testing.T) {
	xtest.SkipUnitIfRequired(t)
	type args struct {
		ctx   context.Context
		table string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Should get constraints",
			args: args{
				ctx:   t.Context(),
				table: tableName,
			},
		},
		{
			name: "Should not get constraints if table not found",
			args: args{
				ctx:   t.Context(),
				table: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedPostgreSQL(t, postgresqlDB, postgresTestDSN)
			t.Cleanup(cleanup)

			e := &postgreSQL{
				db:     db,
				schema: dbName,
			}

			gotRows, _, err := e.GetConstraints(tt.args.ctx, tt.args.table)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, gotRows, "Constraints not found")
		})
	}
}

func newTestPostgreSQL(ctx context.Context) (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("test"),
		postgres.WithUsername("root"),
		postgres.WithPassword("password"),
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("get postgres dsn: %w", err)
	}

	postgresTestDSN = dsn + "sslmode=disable"

	return container, nil
}

func seedPostgreSQL(t *testing.T, driver, dsn string) (*sqlx.DB, func()) {
	ctx := t.Context()
	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, created_at) VALUES
			(1, 'John Doe', 'john@example.com', '2023-01-01 10:00:00'),
			(2, 'Jane Smith', 'jane@example.com', '2023-01-01 10:00:00'),
			(3, 'Bob Wilson', 'bob@example.com', '2023-01-01 10:00:00')
	`)
	require.NoError(t, err)

	return db, func() {
		ctx, cancel := context.WithTimeout(context.Background(), cleanupTimeout)
		defer cancel()

		_, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS users`)
		require.NoError(t, err, "Failed to drop table")
		require.NoError(t, db.Close(), "Failed to close database connection")
	}
}
