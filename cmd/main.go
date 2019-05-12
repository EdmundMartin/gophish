package main

import (
	"fmt"

	"github.com/EdmundMartin/gophish/pkg/engine"
)

func main() {
	piTable := engine.LoadPieceTable("")
	posTable := engine.LoadPositionTable("../configs/positional_tbl.json")
	fmt.Println(posTable)
	finalTables := engine.JoinPosTable(piTable, posTable)
	fmt.Println(finalTables)
	pos := engine.NewGamePosition()
	res := pos.GenerateMoves()
	fmt.Println(res)
	for _, i := range res {
		score := pos.ScoreMove(i, finalTables)
		fmt.Println(i, score)
	}
}
