package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
)

type Parser struct {
	metadata map[string]string

	bilanturiParsers []BilanturiParser
	bilanturiParserForStopCheck BilanturiParser

	repository *Repository

	dbUpdated bool
}

func NewParser(repository *Repository) *Parser {
	if repository == nil {
		repository = NewRepository(config.DBFilePath)
	}

	return &Parser{
		dbUpdated: false,
		repository: repository,
		bilanturiParserForStopCheck: &BilantSimpluParser{},
		bilanturiParsers: []BilanturiParser{
			&BilantSimpluParser{},
			&UUParser{},
			&InstDeCreditParser{},
			&IRParser{},
			&AsigParser{},
			&VSParser{},
			&BrokParser{},
			&VMParser{},
			&IfnParser{},
			&IemeParser{},
			&SifParser{},
			&PensiiParser{},
		},
	}
}

func (this *Parser) expectFieldCount(data []map[string]string, expectedFieldCount int) []map[string]string {
	headerLength := len(data[0]) 
	if headerLength != expectedFieldCount {
		panic(fmt.Errorf("Expected %d fields, got %d", expectedFieldCount, headerLength))
	}

	return data
}

func (this *Parser) parseFirme() []map[string]string {
	data := readData(filepath.Join(config.DataDirectory, "od_firme.csv"))

	return this.expectFieldCount(data, 20)
}

func (this *Parser) getFirmeDataset() string {
	return this.metadata["od_firme.csv"]
}

func (this *Parser) parseReprezentanti() []map[string]string {
	data := readData(filepath.Join(config.DataDirectory, "od_reprezentanti_legali.csv"))
	return this.expectFieldCount(data, 10)
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
	this.expectFieldCount(data, 2)

	nomenclatura := readData(filepath.Join(config.DataDirectory, "n_stare_firma.csv"))
	this.expectFieldCount(nomenclatura, 2)

	this.resolveStariNomenclatura(data, nomenclatura)

	return data
}

func (this *Parser) getStariDataset() string {
	return this.metadata["od_stare_firma.csv"]
}

func (this *Parser) parseCaen() []map[string]string {
	data := readData(filepath.Join(config.DataDirectory, "od_caen_autorizat.csv"))

	return this.expectFieldCount(data, 3)
}

func (this *Parser) getCaenDataset() string {
	return this.metadata["od_caen_autorizat.csv"]
}

func (this *Parser) parseDescriereCaen() []map[string]string {
	data := readData(filepath.Join(config.DataDirectory, "n_caen.csv"))
	return this.expectFieldCount(data, 7)
}

func (this *Parser) getDescriereCaenDataset() string {
	return this.metadata["n_caen.csv"]
}

func (this *Parser) parseDateIdentificare() []map[string]string {
	data := readData(filepath.Join(config.DataDirectory, "od_dateidentificare.txt"))
	return this.expectFieldCount(data, 61)
}

func (this *Parser) getDateIdentificareDataset() string {
	return this.metadata["od_dateidentificare.txt"]
}

func (this *Parser) parseBilanturiForAn(an int) bool {
	anString := strconv.Itoa(an)

	fileName := this.bilanturiParserForStopCheck.GetFileName(anString)
	_, err := os.Stat(filepath.Join(config.DataDirectory, fileName))
	if err != nil {
		return false
	}

	for _, parser := range this.bilanturiParsers {
		group := parser.GetGroupName()
		if this.repository.DoesBilanturiSetExist(an, group) {
			continue
		}

		if !parser.Available(an) {
			continue
		}

		fileName = parser.GetFileName(anString)
		fullPath := filepath.Join(config.DataDirectory, fileName)

		_, err := os.Stat(fullPath)
		if err != nil {
			panic(fmt.Errorf("file %s should exist", fullPath))
		}

		file, err := os.Open(fullPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		logger.Info("reading: ", fullPath)

		data := parser.Parse(file, an)

		fieldCount := len(data[0])
		if !parser.IsFieldCountExpected(an, fieldCount) {
			panic(fmt.Errorf("Unexpected field count: %d, for file: %s", fieldCount, fullPath))
		}

		this.repository.UpdateBilanturi(data, an, group)
		this.setDbUpdated()
	}

	return true
}

func (this *Parser) ParseBilanturi() {
	for i := 2011; ; i++ {
		shouldContinue := this.parseBilanturiForAn(i)
		if !shouldContinue {
			break
		}
	}
}

func (this *Parser) ParseByGrup(reader io.Reader, an int, group string) map[string]int {
	for _, parser := range this.bilanturiParsers {
		if parser.GetGroupName() != group {
			continue
		}

		data := parser.Parse(reader, an)

		if !parser.IsFieldCountExpected(an, len(data[0])) {
			return nil
		}

		return data[0]
	}

	return nil
}

func (this *Parser) setDbUpdated() {
	this.dbUpdated = true
}

func (this *Parser) IsDbUpdated() bool {
	return this.dbUpdated
}

func (this *Parser) Parse() {
	this.repository.Init()

	this.metadata = readMetadata(filepath.Join(config.DataDirectory, "metadata.json"))

	datasetName := this.getFirmeDataset()
	if !this.repository.IsFirmeOnDataset(datasetName) {
		this.repository.UpdateFirme(this.parseFirme(), datasetName)
		this.setDbUpdated()
	}

	datasetName = this.getReprezentantiDataset()
	if !this.repository.IsReprezentantiOnDataset(datasetName) {
		this.repository.UpdateReprezentanti(this.parseReprezentanti(), datasetName)
		this.setDbUpdated()
	}

	datasetName = this.getStariDataset()
	if !this.repository.IsStariOnDataset(datasetName) {
		this.repository.UpdateStari(this.parseStari(), datasetName)
		this.setDbUpdated()
	}

	datasetName = this.getCaenDataset()
	if !this.repository.IsCaenOnDataset(datasetName) {
		this.repository.UpdateCaen(this.parseCaen(), datasetName)
		this.setDbUpdated()
	}

	datasetName = this.getDescriereCaenDataset()
	if !this.repository.IsDescriereCaenOnDataset(datasetName) {
		this.repository.UpdateDescriereCaen(this.parseDescriereCaen(), datasetName)
		this.setDbUpdated()
	}

	datasetName = this.getDateIdentificareDataset()
	if !this.repository.IsDateIdentificareOnDataset(datasetName) {
		this.repository.UpdateDateIdentificare(this.parseDateIdentificare(), datasetName)
		this.setDbUpdated()
	}

	this.ParseBilanturi()

	if this.IsDbUpdated() {
		this.repository.DoOptimizations()
	}

	logger.Info("Finished")
}

