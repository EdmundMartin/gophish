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
	pos := engine.NewGamePosition()
	res := pos.BestMove(finalTables)
	fmt.Println(res)
	pos.Move(res)
}
