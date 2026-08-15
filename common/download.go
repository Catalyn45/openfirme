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

func findDatasets() (string, string){
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

	err := c.Visit("https://data.gov.ro/organization/onrc")
	if err != nil {
		panic(err)
	}

	return firme, nomenclatoare
}

func findResources(url string) []string {
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
			resources = append(resources, e.Attr("href"))
		}
	})

	err := c.Visit(url)
	if err != nil {
		panic(err)
	}

	return resources
}

func downloadFile(url string) (err error) {
	fileName := path.Base(url)
	filePath := "./data/" + fileName

	// Create the file
	out, err := os.Create(filePath)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		return err
	}

	return nil
}

func DownloadResources() {
	firme, nomenclatoare := findDatasets()

	resources := findResources("https://data.gov.ro" + firme)
	resources = append(resources, findResources("https://data.gov.ro" + nomenclatoare)...)

	for _, resource := range resources {
		go downloadFile(resource)
	}
}
