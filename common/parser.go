package common

import (
	"fmt"
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

func (this *Parser) parseBilantSimplu(an int) ([]map[string]int, bool) {
	filePath := filepath.Join(this.dataDirectory, "web_bl_bs_sl_an" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{ 13 }
	if an <= 2015 {
		skipIndexes = append(skipIndexes, 14)
	}

	parsed := readDataDelimiter(filePath, ",", skipIndexes)

	return convertValuesToInt(parsed), true
}

func (this *Parser) parseUU(an int) ([]map[string]int, bool) {
	baseFileName := "web_uu_"
	if an > 2012 {
		baseFileName += "an"
	}

	filePath := filepath.Join(this.dataDirectory, baseFileName + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{ }
	if an <= 2015 {
		skipIndexes = append(skipIndexes, 13)
	}

	parsed := readDataDelimiter(filePath, ",", skipIndexes)

	return convertValuesToInt(parsed), true
}

func (this *Parser) parseIR(an int) ([]map[string]int, bool) {
	// there is no data for 2011
	if an == 2011 {
		return []map[string]int{}, true
	}

	filePath := filepath.Join(this.dataDirectory, "web_ir_an" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{}
	if an <= 2015 && an >= 2018 {
		skipIndexes = append(skipIndexes, 13)
	}

	parsed := readDataDelimiter(filePath, ",", skipIndexes)

	return convertValuesToInt(parsed), true
}

func (this *Parser) parseSituatiiFinanciare(an int) ([]map[string]int, bool) {
	result := []map[string]int{}

	blsl, found := this.parseBilantSimplu(an)
	if !found {
		return nil, false
	}
	result = append(result, blsl...)

	uu, found := this.parseUU(an)
	if !found {
		panic(fmt.Errorf("uu file should exist"))
	}
	result = append(result, uu...)

	ir, found := this.parseIR(an)
	if !found {
		panic(fmt.Errorf("ir file should exist"))
	}
	result = append(result, ir...)

	return result, true
}

func (this *Parser) parseDateIdentificare() []map[string]string {
	return readData(filepath.Join(this.dataDirectory, "od_dateidentificare.txt"))
}

func (this *Parser) getDateIdentificareDataset() string {
	return this.metadata["od_dateidentificare.txt"]
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

	datasetName = this.getDateIdentificareDataset()
	if !this.repository.IsDateIdentificareOnDataset(datasetName) {
		this.repository.UpdateDateIdentificare(this.parseDateIdentificare(), datasetName)
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
