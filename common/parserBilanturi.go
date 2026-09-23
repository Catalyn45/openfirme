package common

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
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

	if an >= 2025 {
		this.expectFieldCount(parsed, 21)
	}

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

	if an >= 2025 {
		this.expectFieldCount(parsed, 21)
	}

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

	if an >= 2025 {
		this.expectFieldCount(parsed, 25)
	}

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

	if an >= 2025 {
		this.expectFieldCount(parsed, 21)
	}

	return convertValuesToInt(parsed), true
}

func (this *Parser) normalizeAsig(data []map[string]int, an int) []map[string]int {
	for _, item := range data {
		if an < 2024 {
			item["I1"] = item["I1"] + item["I3"] + item["I4"]
			item["I2"] =  item["I2"] - item["I3"] - item["I4"] +
							item["I5"] + item["I6"] + item["I7"] +
							item["I8"] + item["I9"] + item["I10"] + item["I11"] + item["I12"]
			item["I7"] = item["I20"]
			item["I10"] = item["I13"]
			item["I17"] = item["I31"]
			item["I18"] = item["I32"]
			item["I19"] = item["I33"]

		} else {
			item["I1"] = item["I1"] + item["I3"]
			item["I2"] =  item["I2"] + item["I8"] + item["I9"] + item["I10"] + item["I14"]
			item["I7"] = item["I24"]
			item["I10"] = item["I15"]
			item["I17"] = item["I36"]
			item["I18"] = item["I37"]
			item["I19"] = item["I38"]
		}

		for key, _ := range item {
			if !slices.Contains([]string{"CUI", "CAEN", "I1", "I2", "I7", "I10", "I17", "I18", "I19"}, key) {
				item[key] = 0
			}
		}
	}

	return data
}

func (this *Parser) parseAsig(an int) ([]map[string]int, bool) {
	filePath := filepath.Join(config.DataDirectory, "webasig" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	parsed := readDataDelimiter(filePath, ",", []int{})
	if an >= 2025 {
		this.expectFieldCount(parsed, 40)
	}

	converted := convertValuesToInt(parsed)

	return this.normalizeAsig(converted, an), true
}

func (this *Parser) normalizeVs(data []map[string]int) []map[string]int {
	for _, item := range data {
		item["I5"] = item["I4"]
		item["I4"] = item["I3"]
		item["I3"] = 0
		item["I7"] = item["I7"] + item["I8"]
		item["I8"] = item["I9"]
		item["I9"] = item["I10"]
		item["I10"] = item["I11"]
		item["I11"] = item["I12"]
		item["I12"] = item["I18"]
		item["I13"] = item["I19"]
		item["I14"] = item["I22"]
		item["I15"] = item["I25"]
		item["I16"] = item["I26"]
		item["I17"] = item["I27"]
		item["I18"] = item["I28"]
		item["I19"] = item["I29"]
	}

	return data
}

func (this *Parser) parseVs(an int) ([]map[string]int, bool) {
	if an < 2015 {
		return []map[string]int{}, true
	}

	filePath := filepath.Join(config.DataDirectory, "web_vs_" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{}
	if an >= 2019 {
		skipIndexes = []int{18, 19}
	}

	parsed := readDataDelimiter(filePath, ",", skipIndexes)
	if an >= 2025 {
		this.expectFieldCount(parsed, 31)
	}

	converted := convertValuesToInt(parsed)

	return this.normalizeVs(converted), true
}

func (this *Parser) normalizeBrok(data []map[string]int, an int) []map[string]int {
	for _, item := range data {
		if an >= 2024 {
			item["I9"], item["I10"] = item["I10"], item["I9"]
		}

		item["I2"] = item["I5"]
		item["I3"] = 0
		item["I4"] = item["I7"]
		item["I5"] = item["I6"]
		item["I6"] = item["I8"]
		item["I7"] = item["I10"]
		item["I8"] = 0
		item["I10"] = item["I11"]
		item["I11"] = item["I12"]
		item["I12"] = item["I15"]
		item["I13"] = item["I18"]
		item["I14"] = item["I19"]
		item["I15"] = item["I20"]
		item["I16"] = item["I21"]
		item["I17"] = item["I22"]
		item["I18"] = item["I23"]
		item["I19"] = item["I24"]
	}

	return data
}

func (this *Parser) parseBrok(an int) ([]map[string]int, bool) {
	filePath := filepath.Join(config.DataDirectory, "webbrok" + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	skipIndexes := []int{}
	if an > 2023 {
		skipIndexes = []int{12, 18}
	}
	parsed := readDataDelimiter(filePath, ",", skipIndexes)

	if an >= 2025 {
		this.expectFieldCount(parsed, 26)
	}

	converted := convertValuesToInt(parsed)

	return this.normalizeBrok(converted, an), true
}

func (this *Parser) normalizeVm(data []map[string]int) []map[string]int {
	for _, item := range data {
		item["I5"] = item["I4"]
		item["I4"] = item["I3"]
		item["I3"] = 0

		item["I7"] = item["I7"] + item["I8"]
		item["I8"] = item["I9"]
		item["I9"] = item["I10"]
		item["I10"] = item["I11"]
		item["I11"] = item["I12"]
		item["I12"] = item["I16"]
		item["I13"] = item["I17"]
		item["I14"] = item["I20"]
		item["I15"] = item["I23"]
		item["I16"] = item["I24"]
		item["I17"] = item["I25"]
		item["I18"] = item["I26"]
		item["I19"] = item["I27"]
	}

	return data
}

func (this *Parser) parseVm(an int) ([]map[string]int, bool) {
	if an < 2016 {
		return []map[string]int{}, true
	}

	baseFileName := "web_vm_"
	if an >= 2021 {
		baseFileName += "an"
	}

	filePath := filepath.Join(config.DataDirectory, baseFileName + strconv.Itoa(an) + ".txt")
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, false
	}

	parsed := readDataDelimiter(filePath, ",", []int{})

	converted := convertValuesToInt(parsed)

	return this.normalizeVm(converted), true
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

	asig, found := this.parseAsig(an)
	if !found {
		panic(fmt.Errorf("asig file should exist"))
	}
	result = append(result, asig...)

	vs, found := this.parseVs(an)
	if !found {
		panic(fmt.Errorf("vs file should exist"))
	}
	result = append(result, vs...)

	brok, found := this.parseBrok(an)
	if !found {
		panic(fmt.Errorf("brok file should exist"))
	}
	result = append(result, brok...)

	vm, found := this.parseVm(an)
	if !found {
		panic(fmt.Errorf("vm file should exist"))
	}
	result = append(result, vm...)

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
