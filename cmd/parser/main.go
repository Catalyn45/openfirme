package main

import (
	"fmt"
	"openfirme/common"
)

func main() {
	// common.DownloadResources()

	repository := common.NewRepository()
	repository.Init()
	repository.Update()

	fmt.Println("Finished")
}