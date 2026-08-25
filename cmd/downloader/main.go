package main

import "openfirme/common"

func main() {
	downloader := common.NewDownloader("https://data.gov.ro/api/3/action", "./data")
	downloader.DownloadData()
}
