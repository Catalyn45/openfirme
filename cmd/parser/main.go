package main

import (
	"openfirme/common"
	"os"
)

func main() {
	parser := common.NewParser(nil)
	parser.Parse()

	if !parser.IsDbUpdated() {
		os.Exit(1)
	}
}