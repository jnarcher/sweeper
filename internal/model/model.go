package model

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jnarcher/sweeper/internal/board"
	"github.com/jnarcher/sweeper/internal/theme"
)

type CursorMotion int

const (
	CursorUp CursorMotion = iota
	CursorDown
	CursorLeft
	CursorRight
)

type GameState int

const (
	Playing GameState = iota
	Win
	Lost
)

type Model struct {
	winWidth  int
	winHeight int
	board     board.Board
	cursor    int
	colors    [9]lipgloss.Color
	State     GameState
}

func InitialModel(config board.BoardConfig) Model {
	cursor := 0
	board := board.NewBoard(config)
	return Model{
		100,
		20,
		board,
		cursor,
		theme.NumberColors,
		Playing,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

			// cursor movement
		case "w", "k", "up":
			if m.State == Playing {
				m.MoveCursor(CursorUp)
			}
			return m, nil
		case "s", "j", "down":
			if m.State == Playing {
				m.MoveCursor(CursorDown)
			}
			return m, nil
		case "a", "h", "left":
			if m.State == Playing {
				m.MoveCursor(CursorLeft)
			}
			return m, nil
		case "d", "l", "right":
			if m.State == Playing {
				m.MoveCursor(CursorRight)
			}
			return m, nil

		// actions
		case " ":
			isBomb := m.board.RevealSquare(m.cursor)

			if isBomb {
				m.SetState(Lost)
                m.board.RevealAll()
				return m, nil
			}

			return m, nil
		case "enter":
			m.board.ToggleFlag(m.cursor)
			if m.CheckWin() {
				m.SetState(Win)
                m.board.RevealAll()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	// draw header
	header := fmt.Sprintf("%d/%d\n", len(m.board.Flagged), len(m.board.Bombs))

	// draw board
	board := ""
	for i := 0; i < m.board.Height; i++ {
		for j := 0; j < m.board.Width; j++ {
			index := i*m.board.Width + j

			if j > 0 {
				board += " "
			}

            style := lipgloss.NewStyle()
			if index == m.cursor {
				style = theme.CursorText
			} else if m.board.Revealed[index] {

                if m.board.IsBomb(index) {
                    style = theme.DefaultText.Foreground(lipgloss.Color("#FF00000"))
                } else {
                    style = theme.DefaultText.Foreground(m.colors[m.board.Squares[index]])
                }
			}

			board += style.Render(m.board.RenderSquare(index))
		}

		if i < m.board.Height-1 {
			board += "\n"
		}
	}
	board = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Render(board)

	// draw win message
	footer := "\n"
	switch m.State {
	case Playing:
		break
	case Win:
		footer += "YOU WIN    Press \"q\" or \"ctrl+q\" to quit"
		break
	case Lost:
		footer += "YOU LOSE   Press \"q\" or \"ctrl+q\" to quit"
		break
	}

	return header + board + footer
}

func (m *Model) MoveCursor(dir CursorMotion) {
	row := m.cursor / m.board.Width
	col := m.cursor % m.board.Width

	switch dir {

	case CursorUp:
		if row > 0 {
			m.cursor -= m.board.Width
		}
		break

	case CursorDown:
		if row < m.board.Height-1 {
			m.cursor += m.board.Width
		}
		break

	case CursorLeft:
		if col > 0 {
			m.cursor--
		}
		break

	case CursorRight:
		if col < m.board.Width-1 {
			m.cursor++
		}
		break
	}
}

func (m Model) CheckWin() bool {
	for _, bomb := range m.board.Bombs {
		if !m.board.IsFlagged(bomb) {
			return false
		}
	}
	return true
}

func (m *Model) SetState(state GameState) {
	m.State = state
	if state != Playing {
		m.cursor = -1
	} else {
		m.cursor = 0
	}
}
