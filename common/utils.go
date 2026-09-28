package common

import (
	"encoding/json"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var allFormeJuridice = []string{
	"AF",
	"ALT",
	"CA",
	"GEIE",
	"GIE",
	"IF",
	"II",
	"INCD",
	"N/A",
	"OC",
	"OC1",
	"OC2",
	"OC3",
	"OC4",
	"OC5",
	"OC6",
	"OC7",
	"OCC",
	"OCM",
	"OCR",
	"PF",
	"PFA",
	"RA",
	"SA",
	"SC",
	"SCA",
	"SCE",
	"SCS",
	"SE",
	"SNC",
	"SRL",
}

func readMetadata(filePath string) map[string]string {
	metadataPath := filePath

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

func saveMetadata(obj map[string]string, filepath string) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		panic(err)
	}
}

func convertValuesToInt(oldmaps []map[string]string) []map[string]int {
	newMaps := []map[string]int{}

	for _, m := range oldmaps {
		newMap := make(map[string]int)
		for key, value := range m {
			newValue := 0

			if value != "" {
				var err error
				newValue, err = strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
			}

			newMap[key] = newValue
		}

		newMaps = append(newMaps, newMap)
	}

	return newMaps
}

func convertDate(datetime string) string {
	if datetime == "" {
		return ""
	}

	layouts := []string{
		"02/01/2006",
		"02/01/2006 15:04",
		"02/01/2006 15:04:05",
	}

	var err error
	for _, layout := range layouts {
		var t time.Time
		t, err = time.Parse(layout, datetime)
		if err == nil {
			return t.Format("2006-01-02 15:04")
		}
	}

	panic(err)
}

func normalize(s string) string {
	s = norm.NFD.String(s)

	var b strings.Builder
	for _, r := range s {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}

func restart(routine func()) {
	r := recover()

	if r != nil {
		logger.Critical("worker panic: ", r, ", stack: ", string(debug.Stack()),  "; restarting")

		time.Sleep(2 * time.Second)
		go routine()
	}
}

func withRestart(routine func()) {
	defer restart(routine)

	routine()
}
