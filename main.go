package main

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/logger"
	"github.com/hrvadl/gowatchsql/internal/platform/db"
	"github.com/hrvadl/gowatchsql/internal/service/engine"
	"github.com/hrvadl/gowatchsql/internal/ui/models/welcome"
)

func main() {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}

	l := logger.New(f)
	pool := db.NewPool()

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
