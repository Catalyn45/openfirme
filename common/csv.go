package common

import (
	"bufio"
	"io"
	"os"
	"slices"
	"strings"
)

func readCsv(reader io.Reader, delimiter string, skipIndexes []string) [][]string {
	scanner := bufio.NewScanner(reader)

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
			if lineNumber > 0 && slices.Contains(skipIndexes, data[0][index]) {
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

func readDataDelimiter(reader io.Reader, delimiter string, skipIndexes []string) []map[string]string {

	data := readCsv(reader, delimiter, skipIndexes)

	parsed := parseCsv(data)

	return parsed
}

func readData(filePath string) []map[string]string {
	logger.Info("reading: ", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return readDataDelimiter(file, "^", []string{})
}

