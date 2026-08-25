package common

import (
	"bufio"
	"os"
	"strings"
)

func read_csv(file_path string, delimiter string, skipIndex int) [][]string {
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
			if index == skipIndex && lineNumber > 0 {
				continue
			}

			el = strings.TrimPrefix(el, "\uFEFF")
			normalized = append(normalized, strings.TrimSpace(el))
		}

		data = append(data, normalized)
	}
	
	return data
}

func parse_csv(data [][]string) []map[string]string {
	parsed := []map[string]string{}

	for i := 1; i < len(data); i++ {
		m := make(map[string]string)

		for j, val := range data[i] {
			m[data[0][j]] = val
		}

		parsed = append(parsed, m)
	}

	return parsed
}

func read_data_delimiter(file_path string, delimiter string, skipIndex int) []map[string]string {
	data := read_csv(file_path, delimiter, skipIndex)

	parsed := parse_csv(data)

	return parsed
}

func read_data(file_path string) []map[string]string {
	return read_data_delimiter(file_path, "^", -1)
}
