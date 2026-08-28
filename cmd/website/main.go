package main

import "openfirme/common"

func main() {
	repository := common.NewRepository("./foo.db")
	server := common.NewServer("127.0.0.1", 8080, repository)
	server.Start()
}