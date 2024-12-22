package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/domain/engine"
	"github.com/hrvadl/gowatchsql/internal/platform/cfg"
	"github.com/hrvadl/gowatchsql/internal/platform/db"
	"github.com/hrvadl/gowatchsql/internal/platform/logger"
	"github.com/hrvadl/gowatchsql/internal/ui/models/welcome"
)

const (
	configPathEnvVarName = "XDG_CONFIG_HOME"
	homeEnvVarName       = "HOME"
)

var basePath = cmp.Or(os.Getenv(configPathEnvVarName), os.Getenv(homeEnvVarName))

func main() {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}

	l := logger.New(f)
	cfg, err := cfg.NewFromFile(basePath)
	if err != nil {
		l.Error("Failed to load config", slog.Any("err", err))
		os.Exit(1)
	}

	pool := db.NewPool(cfg)

	defer func() {
		if err := pool.Close(); err != nil {
			l.Error("Failed to close the pool", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	factory := engine.NewFactory(pool)
	p := tea.NewProgram(welcome.NewModel(l, factory))

	slog.SetLogLoggerLevel(slog.LevelDebug)
	l.Info("Starting the program")

	if _, err := p.Run(); err != nil {
		panic(fmt.Sprintf("Alas, there's been an error: %v", err))
	}
}
