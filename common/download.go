package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
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

func (this *Downloader) getJson(url string) map[string]any {
	resp, err := http.Get(this.url + url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("HTTP error: %s", resp.Status))
	}

	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	if data["success"].(bool) != true {
		panic(fmt.Errorf("fail"))
	}

	return data
}

func findId(data []any, filter string) string {
	for _, pkg := range data {
		el := pkg.(map[string]any)
		if strings.Contains(el["name"].(string), filter) {
			fmt.Println("found name: ", el["name"].(string))
			return el["id"].(string)
		}
	}

	return ""
}

func findIds(data []any, filter string) []Dataset {
	results := []Dataset{}
	for _, el := range data {
		pkg := el.(map[string]any)

		name := pkg["name"].(string)
		id := pkg["id"].(string)

		if strings.Contains(name, filter) {
			fmt.Println("found name: ", name)
			results = append(results, Dataset{
				name: name,
				id: id,
			})
		}
	}

	return results
}

func (this *Downloader) findFirmeNomenclaturaDatasets() (string, string) {
	data := this.getJson("/organization_show?include_datasets=true&id=onrc")

	result := data["result"].(map[string]any)

	firme := findId(result["packages"].([]any), "firme-")
	nomenclatoare := findId(result["packages"].([]any), "nomenclatoare")

	return firme, nomenclatoare
}

type Dataset struct {
	name string
	id string
}

func (this *Downloader) findBilanturiDatasets() []Dataset {
	data := this.getJson("/organization_show?include_datasets=true&id=mfp")

	result := data["result"].(map[string]any)

	return findIds(result["packages"].([]any), "situatii_financiare")
}

func (this *Downloader) findResources(id string, filters []string) []string {
	data := this.getJson("/package_show?id=" + id)

	result := data["result"].(map[string]any)

	downloadLinks := []string{}

	outer:
	for _, el := range result["resources"].([]any) {
		resource := el.(map[string]any)

		for _, filter := range filters {
			if !strings.Contains(resource["name"].(string), filter) {
				continue outer
			}
		}

		downloadLinks = append(downloadLinks, resource["datagovro_download_url"].(string))
	}

	fmt.Println("Download links: ", downloadLinks)

	return downloadLinks
}

func (this *Downloader) downloadFile(url string, replaceExisting bool) (err error) {
	fileName := path.Base(url)
	fmt.Println(fileName)
	filePath := this.outputDir + "/" + fileName

	if !replaceExisting {
		_, err := os.Stat(filePath)
		if err == nil {
			// file already exists, skip
			return nil
		}
	}

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
	firme, nomenclatoare := this.findFirmeNomenclaturaDatasets()

	fmt.Println(firme, nomenclatoare)

	bilanturi := this.findBilanturiDatasets()
	fmt.Println(bilanturi)

	firmeResources := this.findResources(firme, nil)
	firmeResources = append(firmeResources, this.findResources(nomenclatoare, nil)...)

	for _, resource := range firmeResources {
		fmt.Println("Downloading file: ", resource)
		this.downloadFile(resource, true)
		fmt.Println("Finished file: ", resource)
	}

	bilanturiResources := []string{}
	for _, bilant := range bilanturi {
		// there is situatii_financiare_2024_actualizat
		if bilant.name == "situatii_financiare_2024" {
			continue
		}

		// dataset 2020, also contains entries for under 2020
		if strings.Contains(bilant.name, "situatii_financiare_201") {
			continue
		}

		if strings.Contains(bilant.name, "situatii_financiare_200") {
			continue
		}

		bilanturiResources = append(
			bilanturiResources,
			this.findResources(bilant.id, []string{"WEB_BL_BS_SL_AN", ".txt"})...
		)
	}

	for _, resource := range bilanturiResources {
		fmt.Println("Downloading file: ", resource)
		this.downloadFile(resource, false)
		fmt.Println("Finished file: ", resource)
	}
}
