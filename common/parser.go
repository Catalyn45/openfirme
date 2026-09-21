package common

import (
	"path/filepath"
	"slices"
)

type Parser struct {
	metadata map[string]string

	repository *Repository
}

func NewParser() *Parser {
	return &Parser{
		repository: NewRepository(config.DBFilePath),
	}
}

func (this *Parser) parseFirme() []map[string]string {
	return readData(filepath.Join(config.DataDirectory, "od_firme.csv"))
}

func (this *Parser) getFirmeDataset() string {
	return this.metadata["od_firme.csv"]
}

func (this *Parser) parseReprezentanti() []map[string]string {
	return readData(filepath.Join(config.DataDirectory, "od_reprezentanti_legali.csv"))
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
	data := readData(filepath.Join(config.DataDirectory, "od_stare_firma.csv"))

	nomenclatura := readData(filepath.Join(config.DataDirectory, "n_stare_firma.csv"))
	this.resolveStariNomenclatura(data, nomenclatura)

	return data
}

func (this *Parser) getStariDataset() string {
	return this.metadata["od_stare_firma.csv"]
}

func (this *Parser) parseCaen() []map[string]string {
	return readData(filepath.Join(config.DataDirectory, "od_caen_autorizat.csv"))
}

func (this *Parser) getCaenDataset() string {
	return this.metadata["od_caen_autorizat.csv"]
}

func (this *Parser) parseDescriereCaen() []map[string]string {
	return readData(filepath.Join(config.DataDirectory, "n_caen.csv"))
}

func (this *Parser) getDescriereCaenDataset() string {
	return this.metadata["n_caen.csv"]
}

func (this *Parser) parseDateIdentificare() []map[string]string {
	return readData(filepath.Join(config.DataDirectory, "od_dateidentificare.txt"))
}

func (this *Parser) getDateIdentificareDataset() string {
	return this.metadata["od_dateidentificare.txt"]
}

func (this *Parser) Parse() {
	this.repository.Init()

	this.metadata = readMetadata(filepath.Join(config.DataDirectory, "metadata.json"))

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

	datasetName = this.getDescriereCaenDataset()
	if !this.repository.IsDescriereCaenOnDataset(datasetName) {
		this.repository.UpdateDescriereCaen(this.parseDescriereCaen(), datasetName)
	}

	datasetName = this.getDateIdentificareDataset()
	if !this.repository.IsDateIdentificareOnDataset(datasetName) {
		this.repository.UpdateDateIdentificare(this.parseDateIdentificare(), datasetName)
	}

	this.ParseBilanturi()
}
