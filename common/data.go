package common

import (
	"bufio"
	"encoding/json"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func readCsv(file_path string, delimiter string, skipIndexes []int) [][]string {
	file, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	data := [][]string{}

	for lineNumber := 0; scanner.Scan(); lineNumber++ {

		line := scanner.Text()
		splitted := strings.Split(line, delimiter)

		empty := true
		for _, ch := range splitted {
			if ch != "" {
				empty = false
				break
			}
		}

		if empty {
			continue
		}

		normalized := []string{}

		for index, el := range splitted {
			// some old datasets have an additional column we don't care about
			// but it will break the expected order so we just skip
			if slices.Contains(skipIndexes, index) && lineNumber > 0 {
				continue
			}

			el = strings.TrimPrefix(el, "\uFEFF")
			normalized = append(normalized, strings.TrimSpace(el))
		}

		data = append(data, normalized)
	}
	
	return data
}

func parseCsv(data [][]string) []map[string]string {
	parsed := []map[string]string{}

	for i := 1; i < len(data); i++ {
		m := make(map[string]string)

		for j, val := range data[i] {
			if j >= len(data[0]) {
				break
			}

			m[data[0][j]] = val
		}

		parsed = append(parsed, m)
	}

	return parsed
}

func readDataDelimiter(file_path string, delimiter string, skipIndexes []int) []map[string]string {
	data := readCsv(file_path, delimiter, skipIndexes)

	parsed := parseCsv(data)

	return parsed
}

func readData(file_path string) []map[string]string {
	return readDataDelimiter(file_path, "^", []int{})
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
