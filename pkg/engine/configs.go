package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

// PositionTable values of pieces at different positions on the board
type PositionTable struct {
	P []int `json:"P"`
	N []int `json:"N"`
	B []int `json:"B"`
	R []int `json:"R"`
	Q []int `json:"Q"`
	K []int `json:"K"`
}

func (p PositionTable) toMap() map[string][]int {
	return map[string][]int{
		"P": p.P,
		"N": p.N,
		"B": p.B,
		"R": p.R,
		"Q": p.Q,
		"K": p.K,
	}
}

// LoadPieceTable loads a piece table with score values of pieces
func LoadPieceTable(loc string) map[string]int {
	// TODO load values from configs file
	return map[string]int{
		"P": 100,
		"N": 280,
		"B": 320,
		"R": 479,
		"Q": 929,
		"K": 60000,
	}
}

// LoadPositionTable loads a position table from JSON config
func LoadPositionTable(loc string) PositionTable {
	var posTable PositionTable
	fp, err := os.Open(loc)
	defer fp.Close()
	if err != nil {
		fmt.Println("issue loading position table")
		panic(err)
	}
	jsonP := json.NewDecoder(fp)
	jsonP.Decode(&posTable)
	return posTable
}
