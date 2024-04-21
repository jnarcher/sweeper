package tests

import (
	"fmt"
	"testing"

	"github.com/jnarcher/sweeper/internal/board"
)

func boardToString(b board.Board) string {

    s := ""

    for i := 0; i < b.Width; i++ {
        for j := 0; j < b.Height; j++ {
            if j != 0 {
                s += " "
            }
            s += fmt.Sprintf("%d", b.Squares[i * b.Width + j])
        }
        s += "\n"
    }

    return s
}

func TestPrintBoard(t *testing.T) {
    b := board.NewBoard(3, 3, 0)
    bStr := boardToString(b)
    expected := `0 0 0
0 0 0
0 0 0
`

    if bStr != expected {
        t.Fatalf("Board of size (3, 3):\n\n%s\n\nEXPECTED:\n\n%s", bStr, expected)
    }
}
