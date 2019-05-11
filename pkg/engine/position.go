package engine

import "fmt"

// A1 left hand bottom corner of board
const A1 = 91

// H1 right hand bottom corner of board
const H1 = 98

// A8 left hand top corner of board
const A8 = 21

// H8 right hand top corner of board
const H8 = 28

const QS_LIMIT = 150
const EVAL_ROUGHNESS = 20

var validSquares = map[string]bool{
	"p": true,
	"n": true,
	"b": true,
	"q": true,
	"k": true,
	".": true,
}

// PossibleMoves that can be made by a particular piece
var PossibleMoves = map[string][]int{
	"P": []int{-10, -20, -11, -9},
}

// CastleRights represents the right to castle by a particular player
type CastleRights struct {
	East bool
	West bool
}

// Position represents the current position of the board
type Position struct {
	Board       *BoardState
	Score       int
	WhiteCastle *CastleRights
	BlackCastle *CastleRights
	EnPassant   string
	KingPassant string
}

// NewGamePosition creates position for new game
func NewGamePosition() *Position {
	p := &Position{}
	p.Board = NewBoard()
	p.WhiteCastle = &CastleRights{true, true}
	p.BlackCastle = &CastleRights{true, true}
	return p
}

// CalculateMateValues return max and min mate values
func CalculateMateValues(pieceTable map[string]int) (int, int) {
	minMate := pieceTable["K"] - 10*pieceTable["Q"]
	maxMate := pieceTable["K"] + 10*pieceTable["Q"]
	return minMate, maxMate
}

func modifyPosSlice(posSlice []int, pieceValue int) []int {
	mSlice := []int{0}
	for _, item := range posSlice {
		newVal := item + pieceValue
		mSlice = append(mSlice, newVal)
	}
	mSlice = append(mSlice, 0)
	return mSlice
}

func padTable(table []int) []int {
	blanks := []int{}
	for i := 0; i < 20; i++ {
		blanks = append(blanks, 0)
	}
	for _, item := range table {
		blanks = append(blanks, item)
	}
	for i := 0; i < 20; i++ {
		blanks = append(blanks, 0)
	}
	return blanks
}

// JoinPosTable takes piece and position configuration and combines them
func JoinPosTable(pieceTable map[string]int, posTable PositionTable) map[string][]int {
	result := make(map[string][]int)
	for key, value := range posTable.toMap() {
		var newTable []int
		fmt.Println(key, value)
		pieceVal := pieceTable[key]
		i := 0
		for {
			if i+8 <= 64 {
				slice := value[i : i+8]
				newSlice := modifyPosSlice(slice, pieceVal)
				for _, i := range newSlice {
					newTable = append(newTable, i)
				}
			} else {
				break
			}
			i += 8
		}
		result[key] = padTable(newTable)
	}
	return result
}

func validSquare(sq string) bool {
	_, yay := validSquares[sq]
	return yay
}

func validPawnMove(mov int, idx int, tgtIdx int, tgtSq string) bool {
	if mov == -20 && idx > A1+10 {
		return false
	}
	// Pawn moving two squares forward cannot take
	if mov == -20 && tgtSq != "." {
		return false
	}
	// Pawn movint one square forward cannot take
	if mov == -10 && tgtSq != "." {
		return false
	}
	if mov == -9 && tgtSq == "." {
		return false
	}
	if mov == -11 && tgtSq == "." {
		return false
	}
	return true
}

// GenerateMoves takes a position and generates the list of possible moves for the engine
func (pos *Position) GenerateMoves() [][2]int {
	moves := [][2]int{}
	for idx, piece := range pos.Board.State {
		movs, ok := PossibleMoves[piece]
		if ok {
			for _, mov := range movs {
				tgtIdx := mov + idx
				tgtSq := pos.Board.State[tgtIdx]
				if validSquare(tgtSq) {
					if piece == "P" && validPawnMove(mov, idx, tgtIdx, tgtSq) {
						moves = append(moves, [2]int{idx, tgtIdx})
					}
				}
			}
		}
	}
	return moves
}
