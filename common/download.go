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

func (this *Downloader) findDatasets(organization string, patterns []string) map[string]string {
	c := colly.NewCollector()

	results := make(map[string]string)

	for _, pattern := range patterns {
		c.OnHTML(pattern, func(e *colly.HTMLElement) {
			_, ok := results[pattern]
			if !ok {
				fmt.Println("Found:", e.Attr("href"))
				results[pattern] = e.Attr("href")
			}
		})
	}

	err := c.Visit(this.url + organization)
	if err != nil {
		panic(err)
	}

	return results
}

func (this *Downloader) findResources(url string, filter string) []string {
	c := colly.NewCollector()

	resources := []string{}
	c.OnHTML(`ul.resource-list`, func(e *colly.HTMLElement) {
		hrefs := e.ChildAttrs("a", "href")
		titles := e.ChildAttrs("a", "title")

		for index, href := range hrefs {
			if !strings.Contains(href, "dataset") {
				continue
			}

			if !strings.Contains(href, "download") {
				continue
			}

			if !strings.Contains(titles[index], filter) {
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
	patterns := []string{`a[href*="/dataset/firme"]`, `a[href*="/dataset/nomenclatoare"]`}

	results  := this.findDatasets("/organization/onrc", patterns)

	firme := results[patterns[0]]
	nomenclatoare := results[patterns[1]]

	patterns = []string{`a[href*="/dataset/situatii_financiare_2025"]`}
	results = this.findDatasets("/organization/mfp", patterns)

	financiare_2025 := results[patterns[0]]

	resources := this.findResources(this.url + firme, "")

	resources = append(resources, this.findResources(this.url + nomenclatoare, "")...)
	resources = append(resources, this.findResources(this.url + financiare_2025, "WEB_BL_BS_SL_AN2025.txt")...)

	for _, resource := range resources {
		fmt.Println("Downloading file: ", resource)
		this.downloadFile(resource)
		fmt.Println("Finished file: ", resource)
	}
}
