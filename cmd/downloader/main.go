package main

import (
	"openfirme/common"
)

func main() {
	downloader := common.NewDownloader()
	downloader.DownloadData()
}
