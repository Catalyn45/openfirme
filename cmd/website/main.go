package main

import "openfirme/common"

func main() {
	generator := common.NewGenerator()
	generator.Generate()

	server := common.NewServer()
	server.Start()
}