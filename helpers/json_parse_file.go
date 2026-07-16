package helpers

import (
	"encoding/json"
	"fmt"
	"os"
)

func ParseJSONFile(filePath string) (map[string]any, error) {
	data, errReadFile := os.ReadFile(filePath)
	if errReadFile != nil {
		return nil,
			fmt.Errorf(
				"failed to read file: %w",
				errReadFile,
			)
	}

	var result map[string]any

	if errUnmarshal := json.Unmarshal(data, &result); errUnmarshal != nil {
		return nil,
			fmt.Errorf(
				"failed to unmarshal JSON: %w",
				errUnmarshal,
			)
	}

	return result, nil
}
