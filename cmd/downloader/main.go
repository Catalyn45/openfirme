package main

import "openfirme/common"

func main() {
	downloader := common.NewDownloader("https://data.gov.ro", "./data")
	downloader.DownloadResources()
}
