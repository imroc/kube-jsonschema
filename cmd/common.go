package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

func writePrettyJson(v interface{}, path string) error {
	fmt.Printf("write %s\n", path)
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
