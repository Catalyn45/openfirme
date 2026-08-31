package common

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Parser struct {
	dataDirectory string
	repository *Repository
	metadata map[string]string
}

func NewParser(dataDirectory string, repository *Repository) *Parser {
	return &Parser{
		dataDirectory: dataDirectory,
		repository: repository,
	}
}

func (this *Parser) parseFirme() []map[string]string {
	return readData(filepath.Join(this.dataDirectory, "od_firme.csv"))
}

func (this *Parser) getFirmeDataset() string {
	return this.metadata["od_firme.csv"]
}

func (this *Parser) parseReprezentanti() []map[string]string {
	return readData(filepath.Join(this.dataDirectory, "od_reprezentanti_legali.csv"))
}

func (this *Parser) getReprezentantiDataset() string {
	return this.metadata["od_reprezentanti_legali.csv"]
}

func (this *Parser) resolveStariNomenclatura(dataset []map[string]string, nomenclatura []map[string]string) {
	for _, data := range dataset {
		index := slices.IndexFunc(nomenclatura, func (element map[string]string) bool { return data["COD"] == element["COD"]})
		data["STATUS"] = nomenclatura[index]["DENUMIRE"]
	}
}


func (this *Parser) parseStari() []map[string]string {
	data := readData(filepath.Join(this.dataDirectory, "od_stare_firma.csv"))

	nomenclatura := readData(filepath.Join(this.dataDirectory, "n_stare_firma.csv"))
	this.resolveStariNomenclatura(data, nomenclatura)

	return data
}

func (this *Parser) getStariDataset() string {
	return this.metadata["od_stare_firma.csv"]
}

func (this *Parser) parseCaen() []map[string]string {
	return readData(filepath.Join(this.dataDirectory, "od_caen_autorizat.csv"))
}

func (this *Parser) getCaenDataset() string {
	return this.metadata["od_caen_autorizat.csv"]
}

func (this *Parser) parseSituatiiFinanciare(an int) ([]map[string]int, bool) {
	filePath := filepath.Join(this.dataDirectory, "web_bl_bs_sl_an" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndex := -1
	if an <= 2015 {
		skipIndex = 14
	}

	parsed := readDataDelimiter(filePath, ",", skipIndex)

	return convertValuesToInt(parsed), true
}

func (this *Parser) Parse() {
	this.metadata = readMetadata(filepath.Join(this.dataDirectory, "metadata.json"))

	datasetName := this.getFirmeDataset()
	if !this.repository.IsFirmeOnDataset(datasetName) {
		this.repository.UpdateFirme(this.parseFirme(), datasetName)
	}

	datasetName = this.getReprezentantiDataset()
	if !this.repository.IsReprezentantiOnDataset(datasetName) {
		this.repository.UpdateReprezentanti(this.parseReprezentanti(), datasetName)
	}

	datasetName = this.getStariDataset()
	if !this.repository.IsStariOnDataset(datasetName) {
		this.repository.UpdateStari(this.parseStari(), datasetName)
	}

	datasetName = this.getCaenDataset()
	if !this.repository.IsCaenOnDataset(datasetName) {
		this.repository.UpdateCaen(this.parseCaen(), datasetName)
	}

	for i := 2011; ; i++ {
		if this.repository.DoesAnExist(i) {
			continue
		}

		data, exists := this.parseSituatiiFinanciare(i)
		if !exists {
			break
		}

		this.repository.UpdateBilanturi(data, i)
	}
}
