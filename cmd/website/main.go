package main

import "openfirme/common"

func main() {
	repository := common.NewRepository("./foo.db")
	server := common.NewServer("localhost", 8080, repository)
	server.Start()
}