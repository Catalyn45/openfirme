package main

import "openfirme/common"

func main() {
	repository := common.NewRepository()
	server := common.NewServer("localhost", 8080, repository)
	server.Start()
}