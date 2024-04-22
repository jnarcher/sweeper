package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jnarcher/sweeper/internal/board"
	"github.com/jnarcher/sweeper/internal/model"
)

func main() {
    width := flag.Int("w", 10, "width in columns of the board")
    height := flag.Int("h", 10, "height in rows of the board")
    bombs := flag.Int("b", -1, "number of bombs")
    flag.Parse()

    if *bombs == -1 {
        *bombs = *width * *height / 6
    }

	config := board.BoardConfig{
		Width: *width,
		Height: *height,
		Bombs: *bombs,
	}

	p := tea.NewProgram(model.InitialModel(config), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
