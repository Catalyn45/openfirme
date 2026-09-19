package main

import (
	"fmt"
	"openfirme/common"
)

func main() {
	parser := common.NewParser()
	parser.Parse()

	fmt.Println("Finished")
}