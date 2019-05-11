package engine

import (
	"strings"
)

// BoardState represents an array of 120 items allowing for quick discovery of out of bound moves
type BoardState struct {
	State [120]string
}

// NewBoard creates a new board
func NewBoard() *BoardState {
	b := &BoardState{}
	initialBoard := "         \n         \n rnbqkbnr\n pppppppp\n ........\n ........\n ........\n ........\n PPPPPPPP\n RNBQKBNR\n         \n         \n"
	state := strings.Split(initialBoard, "")
	for i, item := range state {
		b.State[i] = item
	}
	return b
}
