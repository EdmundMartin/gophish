package engine

import (
	"fmt"
	"strings"
)

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
	"N": []int{-19, -8, 12, 21, 19, 8, -12, -21},
	"B": []int{-9, -11, 9, 11},
	"R": []int{-10, 10, 1, -1},
	"Q": []int{1, -1, 10, -10, -9, 9, -11, 11},
	"K": []int{1, -1, 10, -10, -9, 9, -11, 11},
}

// ChessMove represents a move on a Chess board
type ChessMove struct {
	Start int
	End   int
	Type  string
}

func (c ChessMove) String() string {
	return fmt.Sprintf("move from %d to %d", c.Start, c.End)
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
	EnPassant   int
	KingPassant int
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
	val := true
	if mov == -20 && idx > A1+10 {
		val = false
	}
	// Pawn moving two squares forward cannot take
	if mov == -20 && tgtSq != "." {
		val = false
	}
	// Pawn movint one square forward cannot
	if mov == -10 && tgtSq != "." {
		val = false
	}
	if mov == -9 && tgtSq == "." {
		val = false
	}
	if mov == -11 && tgtSq == "." {
		val = false
	}
	return val
}

func (pos *Position) freeRange(tgtRange []int) bool {
	for _, idx := range tgtRange {
		if pos.Board.State[idx] != "." {
			return false
		}
	}
	return true
}

func checkForCastles(pos *Position) []ChessMove {
	castlingMoves := []ChessMove{}
	if pos.Board.State[A1] == "R" && pos.Board.State[95] == "K" && pos.WhiteCastle.West == true && pos.freeRange([]int{92, 93, 94}) {
		castlingMoves = append(castlingMoves, ChessMove{95, 93, "CAW"})
	}
	if pos.Board.State[H1] == "R" && pos.Board.State[95] == "K" && pos.WhiteCastle.East == true && pos.freeRange([]int{96, 97}) {
		castlingMoves = append(castlingMoves, ChessMove{95, 97, "CAE"})
	}
	return castlingMoves
}

// GenerateMoves takes a position and generates the list of possible moves for the engine
func (pos *Position) GenerateMoves() []ChessMove {
	moves := []ChessMove{}
	for idx, piece := range pos.Board.State {
		movs, ok := PossibleMoves[piece]
		if ok {
			for _, mov := range movs {
				tgtIdx := mov + idx
				tgtSq := pos.Board.State[tgtIdx]
				if validSquare(tgtSq) {
					if piece == "P" && validPawnMove(mov, idx, tgtIdx, tgtSq) == true {
						moves = append(moves, ChessMove{idx, tgtIdx, "ST"})
					} else if piece != "P" {
						moves = append(moves, ChessMove{idx, tgtIdx, "ST"})
					}
				}
			}
		}
	}
	castMoves := checkForCastles(pos)
	for _, mov := range castMoves {
		moves = append(moves, mov)
	}
	return moves
}

func isCapture(tgtSq string) bool {
	enPs := map[string]bool{
		"p": true,
		"n": true,
		"b": true,
		"q": true,
		"k": true,
	}
	_, yay := enPs[tgtSq]
	return yay
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ScoreMove takes a ChessMove and gives it a score
func (pos *Position) ScoreMove(move ChessMove, posMap map[string][]int) int {
	from := move.Start
	to := move.End
	fromSq := pos.Board.State[from]
	toSq := pos.Board.State[to]
	score := posMap[fromSq][from] - posMap[fromSq][to]
	if isCapture(toSq) {
		score += posMap[strings.ToUpper(toSq)][119-to]
	}
	// Castling check detection
	if abs(to-pos.KingPassant) < 2 {
		score += posMap["K"][119-to]
	}
	// Castling
	if fromSq == "K" && abs(from-to) == 2 {
		score += posMap["R"][abs(from+to)/2]
		if to < from {
			score -= posMap["R"][A1]
		} else {
			score -= posMap["R"][H1]
		}
	}
	if fromSq == "P" {
		// Calculate pawn promotion score
		if A8 <= to && to <= H8 {
			score += posMap["Q"][to] - posMap["P"][to]
		}
		if to == pos.EnPassant {
			score += posMap["P"][119-(to+10)]
		}
	}
	return score
}

// BestMove returns the best legal move
func (pos *Position) BestMove(posMap map[string][]int) ChessMove {
	// Is horrible, lets rewrite
	allMoves := pos.GenerateMoves()
	rankedMoves := make(map[ChessMove]int)
	for _, move := range allMoves {
		score := pos.ScoreMove(move, posMap)
		rankedMoves[move] = score
	}
	var topMove ChessMove
	var count = 0
	var topVal int
	for key, val := range rankedMoves {
		if count == 0 {
			topMove = key
			topVal = val
		} else {
			if val > topVal {
				topMove = key
			}
		}
		count++
	}
	return topMove
}

func put(pos *Position, start, end int, piece string) {
	pos.Board.State[end] = piece
	pos.Board.State[start] = `.`
}

func (pos *Position) updateCastleRights(from int, fromPiece string) {
	if from == A1 {
		pos.WhiteCastle.West = false
	}
	if from == H1 {
		pos.WhiteCastle.East = false
	}
	if from == A8 {
		pos.BlackCastle.West = false
	}
	if from == H8 {
		pos.BlackCastle.East = true
	}
	if fromPiece == "K" {
		pos.WhiteCastle.West = false
		pos.WhiteCastle.East = false
	}
}

// Move makes a chess move on a board
func (pos *Position) Move(chMove ChessMove) {
	from := chMove.Start
	fromPiece := pos.Board.State[from]
	toPiece := pos.Board.State[chMove.End]
	put(pos, chMove.Start, chMove.End, fromPiece)
	fmt.Println(pos.Board.State)
	pos.updateCastleRights(from, fromPiece)
	// Sort out rook after castle
	if chMove.Type == "CAW" {
		pos.Board.State[H1] = `.`
		pos.Board.State[96] = `R`
	}
	if chMove.Type == "CAE" {
		pos.Board.State[A1] = `.`
		pos.Board.State[94] = "R"
	}
	if fromPiece == "P" {
		// Promote Pawn to Queen when it survives to other side of board
		if chMove.End >= A8 && chMove.End <= H8 {
			pos.Board.State[chMove.End] = `Q`
		}
		// Handle en-passant
		if chMove.End-chMove.Start == -20 {
			pos.EnPassant = from + -10
		}
		if chMove.End-chMove.Start == -9 && toPiece == `.` || chMove.End-chMove.Start == -11 && toPiece == `.` {
			pos.Board.State[chMove.End+10] = `.`
		}
	}
}
