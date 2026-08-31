package main

import (
	"fmt"
	"openfirme/common"
)

func main() {
	repository := common.NewRepository("./foo.db")
	repository.Init()

	parser := common.NewParser("./data", repository)
	parser.Parse()

	fmt.Println("Finished")
}