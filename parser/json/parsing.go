package json

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dyllew3/PerfMeasurer/parser"
)

// Parse a json file containing performance request
func Parse(filename string) (parser.RequestsFile, error) {
	result := parser.RequestsFile{}
	bytes, err := os.ReadFile(filename)

	if err != nil {
		log.Printf("Unble to read from file %s\n", filename)
		return result, err
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Printf("Unable to parse json in file %s\n encountered error %v", filename, err)
		err = fmt.Errorf("encountered error %w", err)
	}
	return result, err
}
