package main

import (
	"log"
	"openfirme/common"
)

func main() {
	parser := common.NewParser(nil)
	parser.Parse()

	log.Println("Finished")
}