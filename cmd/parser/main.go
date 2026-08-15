package main

import (
	"fmt"
	"openfirme/common"
)

func main() {
	repository := common.NewRepository("./foo.db")
	repository.Init()
	repository.Update()

	fmt.Println("Finished")
}