package main

import (
	"log"
	"openfirme/common"
)

func main() {
	downloader := common.NewDownloader()
	downloader.DownloadData()

	log.Println("Finished")
}
