package cfg

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestNewFromFileWithCorrectConfig(t *testing.T) {
	now := time.Now().UTC()
	tmpDir := t.TempDir()
	base := filepath.Join(tmpDir, dir)
	configPath := filepath.Join(tmpDir, dir, filename)

	require.NoError(t, os.Mkdir(base, os.ModePerm), "Не вдалося створити директорію")
	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, filemode)
	require.NoError(t, err, "Не вдалося відкрити")

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(base), "Не вдалося видалити файл")
	})

	want := Config{
		Connections: map[string]Connection{
			"sqlite3://test.db": {
				Name:       "test connection",
				DSN:        "sqlite3://test.db",
				LastUsedAt: now,
			},
		},
	}

	data, err := yaml.Marshal(Config{Connections: want.Connections})
	require.NoError(t, err, "Не вдалося закодувати конфіг")

	_, err = f.Write(data)
	require.NoError(t, err, "Не вдалося записати дані в файл")
	require.NoError(t, f.Close(), "Не вдалося закрити файл")

	got, err := NewFromFile(tmpDir)

	require.NoError(t, err)
	require.Equal(t, want.Connections, got.Connections)
}

func TestNewFromFileWithIncorrectFilepath(t *testing.T) {
	tmpDir := "/incorrect/path"
	_, err := NewFromFile(tmpDir)
	require.Error(t, err)
}

func TestNewFromFileWithIncorrectConfig(t *testing.T) {
	tmpDir := t.TempDir()
	base := filepath.Join(tmpDir, dir)
	configPath := filepath.Join(tmpDir, dir, filename)

	require.NoError(t, os.Mkdir(base, os.ModePerm), "Не вдалося створити директорію")
	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, filemode)
	require.NoError(t, err, "Не вдалося відкрити")

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(base), "Не вдалося видалити файл")
	})

	_, err = f.Write([]byte("\tconfig\n"))
	require.NoError(t, err, "Не вдалося записати дані в файл")
	require.NoError(t, f.Close(), "Не вдалося закрити файл")

	cfg, err := NewFromFile(tmpDir)
	require.Error(t, err)
	require.Nil(t, cfg)
}

func TestConfigSave(t *testing.T) {
	tmpDir := t.TempDir()

	conn := Connection{
		Name: "test connection",
		DSN:  "sqlite3://test.db",
	}

	cfg, err := NewFromFile(tmpDir)
	require.NoError(t, err)
	cfg.AddConnection(t.Context(), conn.Name, conn.DSN)
	t.Cleanup(func() {
		require.NoError(t, cfg.Close(), "Не вдалося закрити файл")
	})

	cfg2, err := NewFromFile(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, cfg2.Close(), "Не вдалося закрити файл")
	})

	got := cfg2.Connections[conn.DSN]
	require.Equal(t, conn.DSN, got.DSN)
	require.Equal(t, conn.Name, got.Name)
}

func TestConfig_GetConnections(t *testing.T) {
	now := time.Now().UTC()

	var (
		conn2 = Connection{
			Name:       "conn2",
			DSN:        "sqlite3://conn2.db",
			LastUsedAt: now.Add(time.Second * -1),
		}

		conn1 = Connection{
			Name:       "conn1",
			DSN:        "sqlite3://conn1.db",
			LastUsedAt: now,
		}

		want = []Connection{conn1, conn2}
	)

	cfg := Config{
		Connections: map[string]Connection{
			"conn2": conn2,
			"conn1": conn1,
		},
	}

	got := cfg.GetConnections(t.Context())
	require.Equal(t, want, got)
}

func TestConfigAddConnection(t *testing.T) {
	tmpDir := t.TempDir()

	conn := Connection{
		Name: "test connection",
		DSN:  "sqlite3://test.db",
	}

	cfg, err := NewFromFile(tmpDir)
	require.NoError(t, err)
	cfg.AddConnection(t.Context(), conn.Name, conn.DSN)
	t.Cleanup(func() {
		require.NoError(t, cfg.Close(), "Не вдалося закрити файл")
	})

	cfg2, err := NewFromFile(tmpDir)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, cfg2.Close(), "Не вдалося закрити файл")
	})

	got := cfg2.Connections[conn.DSN]
	require.Equal(t, conn.DSN, got.DSN)
	require.Equal(t, conn.Name, got.Name)
}

func TestConfigAddConnectionIncorrectName(t *testing.T) {
	tmpDir := t.TempDir()

	conn := Connection{
		DSN: "sqlite3://test.db",
	}

	cfg, err := NewFromFile(tmpDir)
	require.NoError(t, err)
	err = cfg.AddConnection(t.Context(), conn.Name, conn.DSN)
	require.Error(t, err)
}

func TestConfigAddConnectionIncorrectDSN(t *testing.T) {
	tmpDir := t.TempDir()

	conn := Connection{
		Name: "name",
	}

	cfg, err := NewFromFile(tmpDir)
	require.NoError(t, err)
	err = cfg.AddConnection(t.Context(), conn.Name, conn.DSN)
	require.Error(t, err)
}

func TestConfigAddConnectionClose(t *testing.T) {
	cfg, err := NewFromFile(t.TempDir())
	require.NoError(t, err)
	cfg.Close()
	require.False(t, isFileDescriptorValid(int(cfg.file.Fd())))
}

func isFileDescriptorValid(fd int) bool {
	var stat syscall.Stat_t
	err := syscall.Fstat(fd, &stat)
	return err == nil
}
