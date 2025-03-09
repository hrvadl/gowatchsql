package engine

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func Test_sqlite_Execute(t *testing.T) {
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
			db, cleanup := seedSQLite(t)
			t.Cleanup(cleanup)

			e := &sqlite{
				db:     db,
				dbPath: dbName,
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

func Test_sqlite_GetTables(t *testing.T) {
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
			db, cleanup := seedSQLite(t)
			t.Cleanup(cleanup)

			e := &sqlite{
				db:     db,
				dbPath: tt.args.schema,
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

func Test_sqlite_GetColumns(t *testing.T) {
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
			db, cleanup := seedSQLite(t)
			t.Cleanup(cleanup)

			e := &sqlite{
				db:     db,
				dbPath: dbName,
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

func Test_sqlite_GetRows(t *testing.T) {
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
				{"1", "John Doe", "john@example.com", "2023-01-01 10:00:00 +0000 UTC"},
				{"2", "Jane Smith", "jane@example.com", "2023-01-01 10:00:00 +0000 UTC"},
				{"3", "Bob Wilson", "bob@example.com", "2023-01-01 10:00:00 +0000 UTC"},
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
			db, cleanup := seedSQLite(t)
			t.Cleanup(cleanup)

			e := &sqlite{
				db:     db,
				dbPath: dbName,
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

func Test_sqlite_GetConstraints(t *testing.T) {
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
			db, cleanup := seedSQLite(t)
			t.Cleanup(cleanup)

			e := &sqlite{
				db:     db,
				dbPath: dbName,
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

func seedSQLite(t *testing.T) (*sqlx.DB, func()) {
	dir := t.TempDir()
	dbFilepath := filepath.Join(dir, "test.db")
	f, err := os.Create(dbFilepath)
	require.NoError(t, err, "Failed to create test database file")
	defer func() {
		require.NoError(t, f.Close(), "Failed to close file")
	}()

	ctx := t.Context()
	db, err := sqlx.ConnectContext(ctx, sqliteDB, dbFilepath)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			email TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
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
