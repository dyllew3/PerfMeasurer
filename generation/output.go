package generation

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Generate output
func GenerateOutputFile(results Output, outputFilename string) error {
	var content []byte = []byte{}
	var err error = nil
	if strings.HasSuffix(outputFilename, ".json") {
		content, err = json.Marshal(results)
	} else if strings.HasSuffix(outputFilename, ".yaml") || strings.HasSuffix(outputFilename, ".yml") {
		content, err = yaml.Marshal(results)
	} else {
		log.Printf("Not supported file type for %s so generating it as json file", outputFilename)
		content, err = json.Marshal(results)
	}
	if err != nil {
		return err
	}
	err = os.WriteFile(outputFilename, content, 0666)
	return err
}
