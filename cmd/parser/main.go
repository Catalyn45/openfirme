package main

import (
	"log"
	"openfirme/common"
)

func main() {
	parser := common.NewParser()
	parser.Parse()

	log.Println("Finished")
}