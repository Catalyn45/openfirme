package main

import (
	"fmt"
	"openfirme/common"
)

func main() {
	downloader := common.NewDownloader()
	downloader.DownloadData()

	fmt.Println("Finished")
}
