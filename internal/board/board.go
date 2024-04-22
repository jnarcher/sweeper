package board

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Width    int
	Height   int
	Bombs    []int
	Squares  []int
	Revealed []bool
	Flagged  []int
}

type BoardConfig struct {
	Width  int
	Height int
	Bombs  int
}

func NewBoard(config BoardConfig) Board {
	if config.Bombs > config.Width*config.Height {
		panic("bomb count cannot be greater than board size")
	}

	bombs := setBombs(config)
	squares := setSquares(bombs, config.Width, config.Height)

	revealed := make([]bool, config.Width*config.Height)
	flagged := []int{}

	return Board{
		config.Width,
		config.Height,
		bombs,
		squares,
		revealed,
		flagged,
	}
}

func setBombs(config BoardConfig) []int {
	bombs := []int{}

	// randomly place bombs
	bombOptions := make([]int, config.Width*config.Height)
	for i := 0; i < len(bombOptions); i++ {
		bombOptions[i] = i
	}

	for ; config.Bombs > 0; config.Bombs-- {
		optionIndex := rand.Intn(len(bombOptions))
		newBomb := bombOptions[optionIndex]
		bombOptions = append(bombOptions[:optionIndex], bombOptions[optionIndex+1:]...)
		didInsert := false

		for i := 0; i < len(bombs); i++ {
			if newBomb < bombs[i] {

				if i == 0 {
					bombs = append([]int{newBomb}, bombs...)
				} else {
					bombs = append(bombs[:i+1], bombs[i:]...)
					bombs[i] = newBomb
				}

				didInsert = true
				break
			}
		}

		if !didInsert {
			bombs = append(bombs, newBomb)
		}
	}
	return bombs
}

// / Returns a string visualization of the board
func (board Board) String() string {
	return boardToString(board, false)
}

func (board Board) StringRevealed() string {
	return boardToString(board, true)
}

func boardToString(board Board, showUnrevealed bool) string {
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
	board.Revealed[index] = true
	if board.IsBomb(index) {
		return true
	}
	recursiveReveal(board, index)
	return false
}

func recursiveReveal(board *Board, focus int) {
	board.Revealed[focus] = true

	// stop if square has a neighboring bomb
	if board.Squares[focus] != 0 {
		return
	}

	unrevealedNeighbors := []int{}
	neighbors := getSurroundingIndices(focus, board.Width, board.Height)
	for i := 0; i < len(neighbors); i++ {
		if !board.Revealed[neighbors[i]] && !board.IsFlagged(neighbors[i]) {
			unrevealedNeighbors = append(unrevealedNeighbors, neighbors[i])
		}
	}

	for i := 0; i < len(unrevealedNeighbors); i++ {
		recursiveReveal(board, unrevealedNeighbors[i])
	}
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
		return "F"
	}

	return "▢"
}

func setSquares(bombs []int, width, height int) []int {
	squares := make([]int, width*height)

	for i := 0; i < width*height; i++ {
		neighborIdxs := getSurroundingIndices(i, width, height)
		surroundingBombCount := 0
		for j := 0; j < len(neighborIdxs); j++ {
			for k := 0; k < len(bombs); k++ {
				if bombs[k] == neighborIdxs[j] {
					surroundingBombCount++
					break
				}
			}
		}
		squares[i] = surroundingBombCount
	}

	return squares
}

func getSurroundingIndices(i, width, height int) []int {
	idxs := []int{}

	row := i / width
	col := i % width

	// top row
	if row > 0 {
		if col > 0 {
			idxs = append(idxs, i-width-1)
		}
		idxs = append(idxs, i-width)
		if col < width-1 {
			idxs = append(idxs, i-width+1)
		}
	}

	// middle row
	if col > 0 {
		idxs = append(idxs, i-1)
	}
	if col < width-1 {
		idxs = append(idxs, i+1)
	}

	// bottom row
	if row < height-1 {
		if col > 0 {
			idxs = append(idxs, i+width-1)
		}
		idxs = append(idxs, i+width)
		if col < width-1 {
			idxs = append(idxs, i+width+1)
		}
	}

	return idxs
}
