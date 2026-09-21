package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func (this *Parser) parseBilantSimplu(an int) ([]map[string]int, bool) {
	filePath := filepath.Join(config.DataDirectory, "web_bl_bs_sl_an" + strconv.Itoa(an) + ".txt")
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

	filePath := filepath.Join(config.DataDirectory, baseFileName + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{ }
	if an >= 2016 {
		skipIndexes = append(skipIndexes, 13)
	}

	parsed := readDataDelimiter(filePath, ",", skipIndexes)

	return convertValuesToInt(parsed), true
}

func (this *Parser) normalizeInstitDeCredit(bilanturi []map[string]int) []map[string]int {
	for _, item := range bilanturi {
		item["I1"] = item["I7"] + item["I8"] + item["I9"]
		item["I2"] = item["I2"] + item["I3"] + item["I4"] + item["I5"] + item["I6"]


		item["I7"] = item["I10"] + item["I11"] + item["I12"]

		item["I9"] = item["I13"]
		item["I10"] = item["I14"] + item["I15"] + item["I16"]

		cifraAfaceri, found := item["I23"]
		if found {
			item["I12"] = cifraAfaceri
		} else {
			item["I12"] = 0
		}

		// set the rest to 0 except I17 and I12
		item["I3"], item["I4"], item["I5"], item["I6"] = 0, 0, 0, 0
		item["I8"], item["I11"], item["I13"] = 0, 0, 0
		item["I14"], item["I15"], item["I18"], item["I19"] = 0, 0, 0, 0
	}

	return bilanturi
}

func (this *Parser) parseInstitDeCredit(an int) ([]map[string]int, bool) {
	if an <= 2013 {
		return []map[string]int{}, true
	}

	baseFileName := "web_inst"
	if an >= 2024 {
		baseFileName += "it"
	}

	if an == 2023 {
		baseFileName += "decredit_"
	} else {
		baseFileName += "_de_credit_"
	}

	if an >= 2024 {
		baseFileName += "an"
	}

	filePath := filepath.Join(config.DataDirectory, baseFileName + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{}
	parsed := readDataDelimiter(filePath, ",", skipIndexes)
	converted := convertValuesToInt(parsed)

	return this.normalizeInstitDeCredit(converted), true
}

func (this *Parser) parseIR(an int) ([]map[string]int, bool) {
	// there is no data for 2011
	if an == 2011 {
		return []map[string]int{}, true
	}

	filePath := filepath.Join(config.DataDirectory, "web_ir_an" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{}
	if an <= 2015 || an >= 2018 {
		skipIndexes = append(skipIndexes, 13)
	}

	parsed := readDataDelimiter(filePath, ",", skipIndexes)

	return convertValuesToInt(parsed), true
}

func (this *Parser) parseBilanturiForAn(an int) ([]map[string]int, bool) {
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

	institDeCredit, found := this.parseInstitDeCredit(an)
	if !found {
		panic(fmt.Errorf("institDeCredit file should exist"))
	}
	result = append(result, institDeCredit...)

	return result, true
}

func (this *Parser) ParseBilanturi() {
	for i := 2011; ; i++ {
		if this.repository.DoesAnExist(i) {
			continue
		}

		data, exists := this.parseBilanturiForAn(i)
		if !exists {
			break
		}

		this.repository.UpdateBilanturi(data, i)
	}
}
