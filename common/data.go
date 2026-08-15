package common

import (
	"bufio"
	"os"
	"strings"
)

func read_csv(file_path string, delimiter string) [][]string {
	file, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	data := [][]string{}

	for scanner.Scan() {
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

		for _, el := range splitted {
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

func read_data_delimiter(file_path string, delimiter string) []map[string]string {
	data := read_csv(file_path, delimiter)

	parsed := parse_csv(data)

	return parsed
}

func read_data(file_path string) []map[string]string {
	return read_data_delimiter(file_path, "^")
}
