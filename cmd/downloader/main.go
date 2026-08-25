package main

import (
	"fmt"
	"openfirme/common"
)

func main() {
	downloader := common.NewDownloader("https://data.gov.ro/api/3/action", "./data")
	downloader.DownloadData()

	fmt.Println("Finished")
}
