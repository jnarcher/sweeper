package board

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Width       int
	Height      int
	BombCount   int
	Bombs       []int
	Squares     []int
	Revealed    []bool
	Flagged     []int
	IsGenerated bool
}

type BoardConfig struct {
	Width         int
	Height        int
	Bombs         int
	InitialSquare int
}

func NewBoard(config BoardConfig) Board {
	if config.Bombs > config.Width*config.Height-1 {
		panic("bomb count cannot be greater than board size")
	}

	return Board{
		Width:       config.Width,
		Height:      config.Height,
		BombCount:   config.Bombs,
		Bombs:       []int{},
		Squares:     make([]int, config.Width*config.Height),
		Revealed:    make([]bool, config.Width*config.Height),
		Flagged:     []int{},
		IsGenerated: false,
	}
}

func (board *Board) Generate(initialReveal int) {
	if board.IsGenerated {
		return
	}

	// randomly place bombs
	bombOptions := make([]int, board.Width*board.Height)
	for i := 0; i < len(bombOptions); i++ {
		bombOptions[i] = i
	}
	// remove initial reveal square from bomb options
	bombOptions = append(bombOptions[:initialReveal], bombOptions[initialReveal+1:]...)

	for count := board.BombCount; count > 0; count-- {
		optionIndex := rand.Intn(len(bombOptions))
		newBomb := bombOptions[optionIndex]
		bombOptions = append(bombOptions[:optionIndex], bombOptions[optionIndex+1:]...)
		didInsert := false

		for i := 0; i < len(board.Bombs); i++ {
			if newBomb < board.Bombs[i] {

				if i == 0 {
					board.Bombs = append([]int{newBomb}, board.Bombs...)
				} else {
					board.Bombs = append(board.Bombs[:i+1], board.Bombs[i:]...)
					board.Bombs[i] = newBomb
				}

				didInsert = true
				break
			}
		}

		if !didInsert {
			board.Bombs = append(board.Bombs, newBomb)
		}
	}

	// fill square numbers
    for _, bomb := range board.Bombs {
        neighborIdxs := board.GetSurroundingIndices(bomb)
        for _, neighbor := range neighborIdxs {
            board.Squares[neighbor] += 1
        }
    } 

	board.IsGenerated = true
}

// / Returns a string visualization of the board
func (board Board) String() string {
	return board.toString(false)
}

func (board Board) StringRevealed() string {
	return board.toString(true)
}

func (board Board) toString(showUnrevealed bool) string {
	s := ""
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {

			index := i*board.Width + j

			if j != 0 {
				s += " "
			}

			charToAdd := ""
			if board.Revealed[index] || showUnrevealed {
				if board.IsBomb(index) {
					charToAdd = "X"
				} else {
					charToAdd = fmt.Sprint(board.Squares[index])
				}
			} else if board.IsFlagged(index) {
				charToAdd = "F"
			} else {
				charToAdd = "-"
			}

			s += charToAdd
		}
		s += "\n"
	}
	return s
}

// / Returns a string visualization of the board indices
func (board Board) ToStringPositions() string {
	s := ""
	for i := 0; i < board.Height; i++ {
		for j := 0; j < board.Width; j++ {

			index := i*board.Width + j

			if j != 0 {
				s += " "
			}

			s += fmt.Sprint(index)
		}
		s += "\n"
	}
	return s
}

func (board Board) IsBomb(index int) bool {
	for i := 0; i < len(board.Bombs); i++ {
		if board.Bombs[i] == index {
			return true
		}
	}
	return false
}

func (board *Board) ToggleFlag(index int) {
    // Cannot flag already revealed squares
    if board.Revealed[index] {
        return;
    }

	if board.IsFlagged(index) {
		for i := 0; i < len(board.Flagged); i++ {
			if board.Flagged[i] == index {
				board.Flagged = append(board.Flagged[:i], board.Flagged[i+1:]...)
				break
			}
		}
	} else {
		board.SetFlag(index)
	}
}

func (board Board) IsFlagged(index int) bool {
	for i := 0; i < len(board.Flagged); i++ {
		if board.Flagged[i] == index {
			return true
		}
	}
	return false
}

func (board *Board) SetFlag(index int) {
	if !board.IsFlagged(index) {
		didInsert := false
		for i := 0; i < len(board.Flagged); i++ {
			if index < board.Flagged[i] {

				if i == 0 {
					board.Flagged = append([]int{index}, board.Flagged...)
				} else {
					board.Flagged = append(board.Flagged[:i+1], board.Flagged[i:]...)
					board.Flagged[i] = index
				}

				didInsert = true
				break
			}
		}

		if !didInsert {
			board.Flagged = append(board.Flagged, index)
		}
	}
}

// / Sets a square to revealed as well as expands te revealed area to show irrelevant squares.
// / Returns whether the revealed square was a bomb.
func (board *Board) RevealSquare(index int) bool {

	if board.IsFlagged(index) {
		return false
	}

	if board.IsBomb(index) {
		board.Revealed[index] = true
		return true
	}

	// If selected an already revealed square reveal all non flagged surrounding squares
	if board.Revealed[index] && board.Squares[index] > 0 {
		neighbors := board.getUnrevealedNeighbors(index)
		for _, neighbor := range neighbors {
			if !board.IsFlagged(neighbor) {
				if board.IsBomb(neighbor) {
					return true
				}
				recursiveReveal(board, neighbor)
			}
		}
		return false
	}

	recursiveReveal(board, index)
	return false
}

func recursiveReveal(board *Board, focus int) {
	board.Revealed[focus] = true

	// Stop if square has a neighboring bomb or is flagged.
	// Should only hit the is flagged condition if incorrectly
	// flagged.
	if board.Squares[focus] != 0 || board.IsFlagged(focus) {
		return
	}

	unrevealedNeighbors := board.getUnrevealedNeighbors(focus)
	for i := 0; i < len(unrevealedNeighbors); i++ {
		recursiveReveal(board, unrevealedNeighbors[i])
	}
}

func (board Board) getUnrevealedNeighbors(focus int) []int {
	unrevealedNeighbors := []int{}
	neighbors := board.GetSurroundingIndices(focus)
	for i := 0; i < len(neighbors); i++ {
		if !board.Revealed[neighbors[i]] && !board.IsFlagged(neighbors[i]) {
			unrevealedNeighbors = append(unrevealedNeighbors, neighbors[i])
		}
	}
	return unrevealedNeighbors
}

func (board *Board) RevealAll() {
	for i := range board.Revealed {
		board.Revealed[i] = true
	}
}

func (board *Board) HideSquare(index int) {
	board.Revealed[index] = false
}

func (board Board) RenderSquare(index int) string {
	if board.Revealed[index] {
		if board.IsBomb(index) {
			return "✖"
		} else if board.Squares[index] == 0 {
			return " "
		}
		return fmt.Sprint(board.Squares[index])
	}

	if board.IsFlagged(index) {
		return "▣"
	}

	return "▢"
}

func (board Board) GetSurroundingIndices(focus int) []int {
	idxs := []int{}

	row := focus / board.Width
	col := focus % board.Width

	// top row
	if row > 0 {
		if col > 0 {
			idxs = append(idxs, focus-board.Width-1)
		}
		idxs = append(idxs, focus-board.Width)
		if col < board.Width-1 {
			idxs = append(idxs, focus-board.Width+1)
		}
	}

	// middle row
	if col > 0 {
		idxs = append(idxs, focus-1)
	}
	if col < board.Width-1 {
		idxs = append(idxs, focus+1)
	}

	// bottom row
	if row < board.Height-1 {
		if col > 0 {
			idxs = append(idxs, focus+board.Width-1)
		}
		idxs = append(idxs, focus+board.Width)
		if col < board.Width-1 {
			idxs = append(idxs, focus+board.Width+1)
		}
	}

	return idxs
}
