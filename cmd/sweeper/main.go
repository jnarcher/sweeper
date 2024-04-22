package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jnarcher/sweeper/internal/board"
	"github.com/jnarcher/sweeper/internal/model"
)

func main() {

	config := board.BoardConfig{
		Width:  50,
		Height: 40,
		Bombs:  200,
	}

	p := tea.NewProgram(model.InitialModel(config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
