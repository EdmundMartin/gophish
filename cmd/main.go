package main

import (
	"fmt"

	"github.com/EdmundMartin/gophish/pkg/engine"
)

func main() {
	posTable := engine.LoadPositionTable("../configs/positional_tbl.json")
	fmt.Println(posTable)
}
