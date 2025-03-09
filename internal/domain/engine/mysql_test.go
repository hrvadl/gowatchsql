package engine

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

const (
	dbName    = "test"
	tableName = "users"
)

const cleanupTimeout = time.Second * 10

var mysqlTestDSN string

func Test_mySQL_Execute(t *testing.T) {
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
				query: "CREATE TABLE test_table (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(100))",
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
			db, cleanup := seedMySQL(t, mysqlDB, mysqlTestDSN)
			t.Cleanup(cleanup)

			e := &mySQL{
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

func Test_mySQL_GetTables(t *testing.T) {
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
		{
			name: "Should return error if schema is not found",
			args: args{
				ctx:    t.Context(),
				schema: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, cleanup := seedMySQL(t, mysqlDB, mysqlTestDSN)
			t.Cleanup(cleanup)

			e := &mySQL{
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

func Test_mySQL_GetColumns(t *testing.T) {
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
			db, cleanup := seedMySQL(t, mysqlDB, mysqlTestDSN)
			t.Cleanup(cleanup)

			e := &mySQL{
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

func Test_mySQL_GetRows(t *testing.T) {
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
				{"1", "John Doe", "john@example.com", "2023-01-01 10:00:00"},
				{"2", "Jane Smith", "jane@example.com", "2023-01-01 10:00:00"},
				{"3", "Bob Wilson", "bob@example.com", "2023-01-01 10:00:00"},
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
			db, cleanup := seedMySQL(t, mysqlDB, mysqlTestDSN)
			t.Cleanup(cleanup)

			e := &mySQL{
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

func Test_mySQL_GetIndexes(t *testing.T) {
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
			db, cleanup := seedMySQL(t, mysqlDB, mysqlTestDSN)
			t.Cleanup(cleanup)

			e := &mySQL{
				db:     db,
				schema: dbName,
			}

			gotRows, _, err := e.GetIndexes(tt.args.ctx, tt.args.table)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			var contains bool
			for _, row := range gotRows {
				if slices.Contains(row, "PRIMARY") {
					contains = true
					break
				}
			}

			require.True(t, contains, "Primary key not found")
		})
	}
}

func Test_mySQL_GetConstraints(t *testing.T) {
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
			db, cleanup := seedMySQL(t, mysqlDB, mysqlTestDSN)
			t.Cleanup(cleanup)

			e := &mySQL{
				db:     db,
				schema: dbName,
			}

			gotRows, _, err := e.GetConstraints(tt.args.ctx, tt.args.table)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			var contains bool
			for _, row := range gotRows {
				if slices.Contains(row, "PRIMARY") {
					contains = true
					break
				}
			}

			require.True(t, contains, "Primary key constraint not found")
		})
	}
}

func newTestMySQL(ctx context.Context) (*mysql.MySQLContainer, error) {
	container, err := mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	if err != nil {
		return nil, fmt.Errorf("start mysql container: %w", err)
	}

	mysqlTestDSN, err = container.ConnectionString(ctx)
	if err != nil {
		return nil, fmt.Errorf("get mysql dsn: %w", err)
	}

	return container, nil
}

func seedMySQL(t *testing.T, driver, dsn string) (*sqlx.DB, func()) {
	ctx := t.Context()
	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
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
