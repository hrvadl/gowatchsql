package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hrvadl/gowatchsql/internal/logger"
	"github.com/hrvadl/gowatchsql/internal/ui/models/welcome"
)

func main() {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}

	l := logger.New(f)

	p := tea.NewProgram(welcome.NewModel(l))
	if _, err := p.Run(); err != nil {
		panic(fmt.Sprintf("Alas, there's been an error: %v", err))
	}
}
