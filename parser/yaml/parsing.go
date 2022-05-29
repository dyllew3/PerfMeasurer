package yaml_handler

import (
	"log"
	"os"

	parser "github.com/dyllew3/PerfMeasurer/parser"
	"gopkg.in/yaml.v3"
)

// Parse a yaml file
func Parse(filename string) (parser.RequestsFile, error) {
	result := parser.RequestsFile{}
	bytes, err := os.ReadFile(filename)

	if err != nil {
		log.Printf("Unble to parse filename %s\n", filename)
		return result, err
	}
	err = yaml.Unmarshal(bytes, &result)
	return result, err
}
