package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
)

type Downloader struct {
	url string
	outputDir string
	metadata map[string]string
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

func (this *Downloader) getMetadata() map[string]string {
	metadataPath := this.outputDir + "/metadata.json"

	obj := make(map[string]string)

	_, err := os.Stat(metadataPath)
	if err != nil {
		return obj
	}

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		panic(err)
	}


	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}

	return obj
}

func (this *Downloader) saveMetadata(obj map[string]string) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}

	metadataPath := this.outputDir + "/metadata.json"

	err = os.WriteFile(metadataPath, data, 0644)
	if err != nil {
		panic(err)
	}
}

func findDataset(data []any, filter string) *Dataset {
	for _, el := range data {
		pkg := el.(map[string]any)

		name := pkg["name"].(string)
		id := pkg["id"].(string)

		if strings.Contains(name, filter) {
			fmt.Println("found name: ", name)
			return &Dataset{
				name: name,
				id: id,
			}
		}
	}

	return nil
}

func findDatasets(data []any, filter string) []Dataset {
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

func (this *Downloader) findFirmeNomenclaturaDatasets() (*Dataset, *Dataset) {
	data := this.getJson("/organization_show?include_datasets=true&id=onrc")

	result := data["result"].(map[string]any)

	packages := result["packages"].([]any)

	firme := findDataset(packages, "firme-")
	nomenclatoare := findDataset(packages, "nomenclatoare")

	return firme, nomenclatoare
}

type Dataset struct {
	name string
	id string
}

func (this *Downloader) findBilanturiDatasets() []Dataset {
	data := this.getJson("/organization_show?include_datasets=true&id=mfp")

	result := data["result"].(map[string]any)

	return findDatasets(result["packages"].([]any), "situatii_financiare")
}

func (this *Downloader) findResources(id string, filters []string) []string {
	data := this.getJson("/package_show?id=" + id)

	result := data["result"].(map[string]any)

	resources := result["resources"].([]any)

	downloadLinks := []string{}
	actualizatDownloadlinks := []string{}

	outer:
	for _, el := range resources {
		resource := el.(map[string]any)

		resourceName := resource["name"].(string)

		// files with -actualizat are missing extension
		isActualizat := strings.Contains(resourceName, " - actualizat")
		if isActualizat {
			resourceName += ".txt"
		}

		for _, filter := range filters {
			if !strings.Contains(resourceName, filter) {
				continue outer
			}
		}

		downloadUrl := resource["datagovro_download_url"].(string)

		if isActualizat {
			actualizatDownloadlinks = append(actualizatDownloadlinks, downloadUrl)
		} else {
			downloadLinks = append(downloadLinks, downloadUrl)
		}
	}

	if len(actualizatDownloadlinks) > 0 {
		downloadLinks = actualizatDownloadlinks
	}

	fmt.Println("Download links: ", downloadLinks)

	return downloadLinks
}

func (this *Downloader) downloadFile(url string, fileName string) {
	filePath := this.outputDir + "/" + fileName

	out, err := os.Create(filePath)
	if err != nil  {
		panic(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("bad status: %s", resp.Status))
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil  {
		panic(err)
	}
}

func (this *Downloader) downloadResources(dataset *Dataset, resources []string, addOnly bool) {
	for _, resource := range resources {
		fileName := path.Base(resource)

		val, ok := this.metadata[fileName]

		if addOnly && ok {
			continue
		}

		if ok && val == dataset.name {
			// skip as we already have the latest data
			continue
		}

		fmt.Println("Downloading file: ", resource)
		this.downloadFile(resource, fileName)
		fmt.Println("Finished file: ", resource)

		this.metadata[fileName] = dataset.name
		this.saveMetadata(this.metadata)
	}
}

func (this *Downloader) DownloadData() {
	this.metadata = this.getMetadata()

	firme, nomenclatoare := this.findFirmeNomenclaturaDatasets()
	fmt.Println(firme, nomenclatoare)

	bilanturi := this.findBilanturiDatasets()
	fmt.Println(bilanturi)

	firmeResources := this.findResources(firme.id, nil)
	nomenclatoareResources := this.findResources(nomenclatoare.id, nil)

	os.Mkdir(this.outputDir, 0755)

	this.downloadResources(firme, firmeResources, false)
	this.downloadResources(nomenclatoare, nomenclatoareResources, false)

	sort.Slice(bilanturi, func(i, j int) bool {
		return bilanturi[i].name > bilanturi[j].name
	})

	fmt.Println("sorted: ", bilanturi)

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

		bilanturiResources := this.findResources(bilant.id, []string{"WEB_BL_BS_SL_AN", ".txt"})
		this.downloadResources(&bilant, bilanturiResources, true)
	}
}
