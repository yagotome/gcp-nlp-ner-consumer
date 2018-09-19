package csv

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/yagotome/gcp-nlp-ner-consumer/utils/sliceutil"
)

// ExtractColumnGroupedBy ...
func ExtractColumnGroupedBy(filename, columnName, groupColumn string, groupFilter []string) (map[string][]string, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fileContent := string(fileBytes)
	lines := strings.Split(fileContent, "\n")
	headers := strings.Split(lines[0], ";")
	colIdx := sliceutil.IndexOf(headers, columnName)
	if colIdx < 0 {
		return nil, fmt.Errorf("Column %v not found", columnName)
	}
	groupColIdx := sliceutil.IndexOf(headers, groupColumn)
	if groupColIdx < 0 {
		return nil, fmt.Errorf("Column %v not found", groupColumn)
	}
	valueLines := lines[1:]
	valuesGrouped := make(map[string][]string)
	for _, line := range valueLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		columns := strings.Split(line, ";")
		if groupColIdx >= len(columns) || colIdx >= len(columns) {
			return nil, fmt.Errorf("Invalid CSV")
		}
		groupName := columns[groupColIdx]
		if !sliceutil.Contains(groupFilter, groupName) {
			continue
		}
		if _, ok := valuesGrouped[groupName]; !ok {
			valuesGrouped[groupName] = make([]string, 0)
		}
		valuesGrouped[groupName] = append(valuesGrouped[groupName], columns[colIdx])
	}
	return valuesGrouped, nil
}
