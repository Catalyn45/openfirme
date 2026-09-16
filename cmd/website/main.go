package main

import "openfirme/common"

func main() {
	generator := common.NewGenerator("./templates", "./public")
	generator.Generate()

	repository := common.NewRepository("file:foo.db?mode=ro")
	server := common.NewServer("127.0.0.1", 8080, repository)
	server.Start()
}