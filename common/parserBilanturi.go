package common

import (
	"io"
	"slices"
)

type BilanturiParser interface {
	Available(an int) bool
	GetFileName(an string) string
	GetGroupName() string
	IsFieldCountExpected(an int, fieldCount int) bool
	Parse(reader io.Reader, an int) []map[string]int
}

type BilantSimpluParser struct {}

func (this *BilantSimpluParser) Available(_ int) bool {
	return true
}

func (this *BilantSimpluParser) GetFileName(an string) string {
	return "web_bl_bs_sl_an" + an + ".txt"
}

func (this *BilantSimpluParser) GetGroupName() string {
	return "bl_sl"
}

func (this *BilantSimpluParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 21
}

func (this *BilantSimpluParser) Parse(reader io.Reader, an int) []map[string]int {
	skipIndexes := []string{ "I12" }
	if an <= 2015 {
		skipIndexes = append(skipIndexes, "I13")
	}

	parsed := readDataDelimiter(reader, ",", skipIndexes)

	return convertValuesToInt(parsed)
}

type UUParser struct {}

func (this *UUParser) Available(an int) bool {
	return true
}

func (this *UUParser) GetFileName(an string) string {
	fileName := "web_uu_"
	if an > "2012" {
		fileName += "an"
	}

	return fileName + an + ".txt"
}

func (this *UUParser) GetGroupName() string {
	return "uu"
}

func (this *UUParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 21
}

func (this *UUParser) Parse(reader io.Reader, an int) []map[string]int {
	skipIndexes := []string{}
	if an >= 2016 {
		skipIndexes = append(skipIndexes, "I12")
	}

	parsed := readDataDelimiter(reader, ",", skipIndexes)

	return convertValuesToInt(parsed)
}

type InstDeCreditParser struct {}

func (this *InstDeCreditParser) Available(an int) bool {
	return an > 2013
}

func (this *InstDeCreditParser) GetFileName(an string) string {
	fileName := "web_inst"
	if an >= "2024" {
		fileName += "it"
	}

	if an == "2023" {
		fileName += "decredit_"
	} else {
		fileName += "_de_credit_"
	}

	if an >= "2024" {
		fileName += "an"
	}

	return fileName + an + ".txt"
}

func (this *InstDeCreditParser) GetGroupName() string {
	return "inst_credit"
}

func (this *InstDeCreditParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 25
}

func (this *InstDeCreditParser) normalize(bilanturi []map[string]int) []map[string]int {
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

func (this *InstDeCreditParser) Parse(reader io.Reader, an int) []map[string]int {
	parsed := readDataDelimiter(reader, ",", []string{})

	converted := convertValuesToInt(parsed)

	return this.normalize(converted)
}

type IRParser struct {}

func (this *IRParser) Available(an int) bool {
	return an > 2011
}

func (this *IRParser) GetFileName(an string) string {
	return "web_ir_an" + an + ".txt"
}

func (this *IRParser) GetGroupName() string {
	return "ir"
}

func (this *IRParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 21
}

func (this *IRParser) Parse(reader io.Reader, an int) []map[string]int {
	skipIndexes := []string{}
	if an <= 2015 || an >= 2018 {
		skipIndexes = append(skipIndexes, "I12")
	}

	parsed := readDataDelimiter(reader, ",", skipIndexes)

	return convertValuesToInt(parsed)
}


type AsigParser struct {}

func (this *AsigParser) Available(_ int) bool {
	return true
}

func (this *AsigParser) GetFileName(an string) string {
	return "webasig" + an + ".txt"
}

func (this *AsigParser) GetGroupName() string {
	return "asig"
}

func (this *AsigParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 40
}

func (this *AsigParser) normalize(data []map[string]int, an int) []map[string]int {
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

func (this *AsigParser) Parse(reader io.Reader, an int) []map[string]int {
	parsed := readDataDelimiter(reader, ",", []string{})

	converted := convertValuesToInt(parsed)

	return this.normalize(converted, an)
}

type VSParser struct {}

func (this *VSParser) Available(an int) bool {
	return an > 2014
}

func (this *VSParser) GetFileName(an string) string {
	return "web_vs_" + an + ".txt"
}

func (this *VSParser) GetGroupName() string {
	return "vs"
}

func (this *VSParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 31
}

func (this *VSParser) normalize(data []map[string]int) []map[string]int {
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

func (this *VSParser) Parse(reader io.Reader, an int) []map[string]int {
	skipIndexes := []string{}
	if an >= 2019 {
		skipIndexes = []string{"I17", "I18"}
	}

	parsed := readDataDelimiter(reader, ",", skipIndexes)

	converted := convertValuesToInt(parsed)

	return this.normalize(converted)
}

type BrokParser struct {}

func (this *BrokParser) Available(_ int) bool {
	return true
}

func (this *BrokParser) GetFileName(an string) string {
	return "webbrok" + an + ".txt"
}

func (this *BrokParser) GetGroupName() string {
	return "brok"
}

func (this *BrokParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 26
}

func (this *BrokParser) normalize(data []map[string]int, an int) []map[string]int {
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

func (this *BrokParser) Parse(reader io.Reader, an int) []map[string]int {
	skipIndexes := []string{}
	if an > 2023 {
		skipIndexes = []string{"I11", "I17"}
	}
	parsed := readDataDelimiter(reader, ",", skipIndexes)

	converted := convertValuesToInt(parsed)

	return this.normalize(converted, an)
}

type VMParser struct {}

func (this *VMParser) Available(an int) bool {
	return an > 2015
}

func (this *VMParser) GetFileName(an string) string {
	fileName := "web_vm_"
	if an >= "2021" {
		fileName += "an"
	}

	return fileName + an + ".txt"
}

func (this *VMParser) GetGroupName() string {
	return "vm"
}

func (this *VMParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 29
}

func (this *VMParser) normalize(data []map[string]int) []map[string]int {
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

func (this *VMParser) Parse(reader io.Reader, an int) []map[string]int {
	parsed := readDataDelimiter(reader, ",", []string{})

	converted := convertValuesToInt(parsed)

	return this.normalize(converted)
}

type IfnParser struct {}

func (this *IfnParser) Available(an int) bool {
	return an > 2011
}

func (this *IfnParser) GetFileName(an string) string {
	fileName := "web_ifn"
	if an == "2014" {
		fileName += "_"
	}

	return fileName + an + ".txt"
}

func (this *IfnParser) GetGroupName() string {
	return "ifn"
}

func (this *IfnParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 25
}

func (this *IfnParser) normalize(an int, data []map[string]int) []map[string]int {
	for _, item := range data {
		if an >= 2023 {
			item["I5"] = item["I1"]
			item["I1"] = item["I7"] + item["I8"] + item["I9"]
			item["I2"] = 0
			item["I3"] = 0
			item["I4"] = 0
			item["I6"] = 0
			item["I7"] = item["I10"] + item["I11"] + item["I12"]
			item["I8"] = 0
			item["I9"] = item["I13"]
			item["I10"] = item["I14"] + item["I15"] + item["I16"]
			item["I11"] = item["I14"]

			var found bool
			item["I12"], found = item["I23"]
			if !found {
				item["I12"] = 0
			}

			item["I13"] = item["I18"]
			item["I14"] = 0
			item["I15"] = item["I19"]
			item["I16"] = 0
			item["I17"] = item["I22"]
			item["I18"] = 0
			item["I19"] = 0
		} else {
			item["I5"] = item["I1"]
			item["I1"] = item["I8"] + item["I9"]
			item["I4"] = item["I2"] + item["I3"]
			item["I2"] = 0
			item["I3"] = 0
			item["I6"] = 0
			item["I7"] = item["I10"] + item["I11"] + item["I12"]
			item["I8"] = 0
			item["I9"] = item["I13"]
			item["I10"] = item["I15"]
			item["I11"] = item["I15"]
			item["I12"] = 0
			item["I13"] = item["I19"]
			item["I14"] = item["I20"]
			item["I15"] = item["I21"]
			item["I16"] = item["I22"]
			item["I17"] = item["I23"]
			item["I18"] = item["I24"]
			item["I19"] = 0
		}
	}

	return data
}

func (this *IfnParser) Parse(reader io.Reader, an int) []map[string]int {
	parsed := readDataDelimiter(reader, ",", []string{})

	converted := convertValuesToInt(parsed)

	return this.normalize(an, converted)
}

type IemeParser struct {}

func (this *IemeParser) Available(an int) bool {
	return an > 2022
}

func (this *IemeParser) GetFileName(an string) string {
	return "web_ip_ieme" + an + ".txt"
}

func (this *IemeParser) GetGroupName() string {
	return "ieme"
}

func (this *IemeParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 26
}

func (this *IemeParser) normalize(data []map[string]int) []map[string]int {
	for _, item := range data {
		item["I4"] = item["I2"] + item["I3"]
		item["I2"] = item["I1"] + item["I2"] + item["I3"] + item["I4"] + item["I5"] + item["I6"] + item["I7"]
		item["I5"] = item["I1"]
		item["I1"] = item["I8"] + item["I9"]
		item["I3"] = 0
		item["I6"] = 0
		item["I7"] = item["I10"] + item["I11"] + item["I12"] + item["I14"]
		item["I8"] = 0
		item["I9"] = item["I13"]
		item["I10"] = item["I15"] + item["I16"] + item["I17"] - item["I18"]
		item["I11"] = item["I15"]
		item["I12"] = 0
		item["I13"] = item["I19"]
		item["I14"] = item["I20"]
		item["I15"] = item["I21"]
		item["I16"] = item["I22"]
		item["I17"] = item["I23"]
		item["I18"] = item["I24"]
		item["I19"] = 0
	}

	return data
}

func (this *IemeParser) Parse(reader io.Reader, an int) []map[string]int {
	parsed := readDataDelimiter(reader, ",", []string{})

	converted := convertValuesToInt(parsed)

	return this.normalize(converted)
}

type SifParser struct {}

func (this *SifParser) Available(an int) bool {
	return an > 2021
}

func (this *SifParser) GetFileName(an string) string {
	return "web_sif" + an + ".txt"
}

func (this *SifParser) GetGroupName() string {
	return "sif"
}

func (this *SifParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 28
}

func (this *SifParser) normalize(data []map[string]int) []map[string]int {
	for _, item := range data {
		item["I5"] = item["I4"]
		item["I4"] = item["I3"]
		item["I3"] = 0
		item["I7"] = item["I7"] + item["I8"]
		item["I8"] = item["I9"]
		item["I9"] = item["I10"]
		item["I10"] = item["I19"]
		item["I11"] = item["I12"]
		item["I12"] = 0
		item["I13"] = item["I20"]
		item["I14"] = item["I21"]
		item["I15"] = item["I22"]
		item["I16"] = item["I23"]
		item["I17"] = item["I24"]
		item["I18"] = item["I25"]
		item["I19"] = item["I26"]
	}

	return data
}

func (this *SifParser) Parse(reader io.Reader, an int) []map[string]int {
	parsed := readDataDelimiter(reader, ",", []string{})

	converted := convertValuesToInt(parsed)

	return this.normalize(converted)
}

type PensiiParser struct {}

func (this *PensiiParser) Available(_ int) bool {
	return true
}

func (this *PensiiParser) GetFileName(an string) string {
	fileName := "web_pensii"
	if an == "2011" || an == "2012" || an == "2014" {
		fileName += "_"
	}

	return fileName + an + ".txt"
}

func (this *PensiiParser) GetGroupName() string {
	return "pensii"
}

func (this *PensiiParser) IsFieldCountExpected(an int, fieldCount int) bool {
	if an < 2025 {
		return true
	}

	return fieldCount == 21
}


func (this *PensiiParser) Parse(reader io.Reader, an int) []map[string]int {
	skipIndexes := []string{ "I5" }
	parsed := readDataDelimiter(reader, ",", skipIndexes)

	return convertValuesToInt(parsed)
}
