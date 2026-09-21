package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Downloader struct {
	url string
	metadata map[string]string
}

func NewDownloader() *Downloader {
	return &Downloader{
		url: "https://data.gov.ro/api/3/action",
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
	metadataPath := filepath.Join(config.DataDirectory, "metadata.json")
	return readMetadata(metadataPath)
}

func (this *Downloader) saveMetadata(obj map[string]string) {
	metadataPath := filepath.Join(config.DataDirectory, "metadata.json")
	saveMetadata(obj, metadataPath)
}

func findDataset(data []any, filter string) *Dataset {
	for _, el := range data {
		pkg := el.(map[string]any)

		name := pkg["name"].(string)
		id := pkg["id"].(string)

		if strings.Contains(name, filter) {
			log.Println("found name: ", name)
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
			log.Println("found name: ", name)
			results = append(results, Dataset{
				name: name,
				id: id,
			})
		}
	}

	return results
}

func (this *Downloader) getPackages(orgId string) []any {
	data := this.getJson("/organization_show?include_datasets=true&id=" + orgId)

	result := data["result"].(map[string]any)

	return result["packages"].([]any)
}

type Dataset struct {
	name string
	id string
}

func (this *Downloader) findResources(id string, filtersSet [][]string) []string {
	data := this.getJson("/package_show?id=" + id)

	result := data["result"].(map[string]any)

	resources := result["resources"].([]any)

	downloadLinks := []string{}
	actualizatDownloadlinks := []string{}

	for _, el := range resources {
		resource := el.(map[string]any)

		resourceName := strings.ToLower(resource["name"].(string))

		// files with -actualizat are missing extension
		isActualizat := strings.Contains(resourceName, "actualizat")
		if resourceName[len(resourceName)-3:] != "txt" && isActualizat {
			resourceName += ".txt"
		}

		filtersMatched := false

		filtersLabel:
		for _, filters := range filtersSet {
			for _, filter := range filters {
				if !strings.Contains(resourceName, filter) {
					continue filtersLabel
				}
			}

			filtersMatched = true
			break
		}

		if filtersSet != nil && !filtersMatched {
			continue
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

	log.Println("Download links: ", downloadLinks)

	return downloadLinks
}

func (this *Downloader) downloadFile(url string, fileName string, recreate bool) {
	filePath := filepath.Join(config.DataDirectory, fileName)

	flags := os.O_WRONLY|os.O_CREATE
	if recreate {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}

	out, err := os.OpenFile(filePath, flags, 0644)
	if err != nil {
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

	var reader io.Reader = resp.Body
	if !recreate {
		br := bufio.NewReader(resp.Body)

		// Skip first line
		_, err := br.ReadString('\n')
		if err != nil {
			panic(err)
		}

		reader = br
	}

	_, err = io.Copy(out, reader)
	if err != nil  {
		panic(err)
	}
}

func (this *Downloader) downloadResourcesCombined(dataset *Dataset, resources []string, fileName string) {
	val, ok := this.metadata[fileName]
	if ok && val == dataset.name {
		// skip as we already have the latest data
		return
	}

	for index, resource := range resources {
		log.Println("Downloading file: ", resource)
		this.downloadFile(resource, fileName, index == 0)
		log.Println("Finished file: ", resource)
	}

	this.metadata[fileName] = dataset.name
	this.saveMetadata(this.metadata)
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

		log.Println("Downloading file: ", resource)
		this.downloadFile(resource, fileName, true)
		log.Println("Finished file: ", resource)

		this.metadata[fileName] = dataset.name
		this.saveMetadata(this.metadata)
	}
}

func (this *Downloader) DownloadData() {
	this.metadata = this.getMetadata()

	packages := this.getPackages("onrc")

	firme := findDataset(packages, "firme-")
	nomenclatoare := findDataset(packages, "nomenclatoare")

	log.Println(firme, nomenclatoare)

	packages = this.getPackages("mfp")

	dateIdentificare := findDataset(packages, "date_de_identificare_")
	bilanturi := findDatasets(packages, "situatii_financiare")

	log.Println(dateIdentificare, bilanturi)

	firmeResources := this.findResources(firme.id, nil)
	nomenclatoareResources := this.findResources(nomenclatoare.id, nil)
	dateIdentificareResources := this.findResources(dateIdentificare.id, [][]string{
		[]string{ "date_identificare_platitori_", ".txt" },
	})

	os.Mkdir(config.DataDirectory, 0755)

	this.downloadResources(firme, firmeResources, false)
	this.downloadResources(nomenclatoare, nomenclatoareResources, false)
	this.downloadResourcesCombined(dateIdentificare, dateIdentificareResources, "od_dateidentificare.txt")

	sort.Slice(bilanturi, func(i, j int) bool {
		return bilanturi[i].name > bilanturi[j].name
	})

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

		bilanturiResources := this.findResources(bilant.id, [][]string{
			[]string{"web_bl_bs_sl_an", ".txt"},
			[]string{"web_uu_", ".txt"},
			[]string{"web_ir_an", ".txt"},
			[]string{"web_ir_an", ".txt"},
			[]string{"web_inst", "_de_credit_", ".txt"},
			[]string{"web_instdecredit", ".txt"},
			[]string{"webasig", ".txt"},
		})

		this.downloadResources(&bilant, bilanturiResources, true)
	}
}
