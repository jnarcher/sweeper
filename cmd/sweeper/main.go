package main

import (
	"fmt"
	"github.com/jnarcher/sweeper/internal/board"
)

func main() {
    bd := board.NewBoard(5, 5, 5)
    fmt.Println(bd.ToStringRevealed())
    fmt.Println(bd.ToString())

    died := bd.RevealSquare(0)
    if died {
        fmt.Println("Bomb revealed")
    }

    fmt.Println(bd.ToString())
}
