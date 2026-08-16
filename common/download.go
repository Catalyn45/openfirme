package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"path"

	"github.com/gocolly/colly/v2"
)

type Downloader struct {
	url string
	outputDir string
}

func NewDownloader(url string, outputDir string) *Downloader {
	return &Downloader{
		url: url,
		outputDir: outputDir,
	}
}

func (this *Downloader) findDatasets() (string, string){
	c := colly.NewCollector()

	var firme string = ""
	c.OnHTML(`a[href*="/dataset/firme"]`, func(e *colly.HTMLElement) {
		if firme == "" {
			fmt.Println("Found:", e.Attr("href"))
			firme = e.Attr("href")
		}
	})

	var nomenclatoare string = ""
	c.OnHTML(`a[href*="/dataset/nomenclatoare"]`, func(e *colly.HTMLElement) {

		if nomenclatoare == "" {
			fmt.Println("Found:", e.Attr("href"))
			nomenclatoare = e.Attr("href")
		}
	})

	err := c.Visit(this.url + "/organization/onrc")
	if err != nil {
		panic(err)
	}

	return firme, nomenclatoare
}

func (this *Downloader) findResources(url string) []string {
	c := colly.NewCollector()

	resources := []string{}
	c.OnHTML(`ul.resource-list`, func(e *colly.HTMLElement) {
		hrefs := e.ChildAttrs("a", "href")

		for _, href := range hrefs {
			if !strings.Contains(href, "dataset") {
				continue
			}

			if !strings.Contains(href, "download") {
				continue
			}

			fmt.Println("Found:", href)
			resources = append(resources, href)
		}
	})

	err := c.Visit(url)
	if err != nil {
		panic(err)
	}

	fmt.Println("Resources: ", resources)

	return resources
}

func (this *Downloader) downloadFile(url string) (err error) {
	fileName := path.Base(url)
	fmt.Println(fileName)
	filePath := this.outputDir + "/" + fileName

	out, err := os.Create(filePath)
	if err != nil  {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

func (this *Downloader) DownloadResources() {
	firme, nomenclatoare := this.findDatasets()

	resources := this.findResources(this.url + firme)
	resources = append(resources, this.findResources(this.url + nomenclatoare)...)

	for _, resource := range resources {
		fmt.Println("Downloading file: ", resource)
		this.downloadFile(resource)
		fmt.Println("Finished file: ", resource)
	}
}
