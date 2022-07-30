package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dyllew3/PerfMeasurer/generation"
	"github.com/dyllew3/PerfMeasurer/parser"
	json_handler "github.com/dyllew3/PerfMeasurer/parser/json"
	yaml_handler "github.com/dyllew3/PerfMeasurer/parser/yaml"
	"github.com/dyllew3/PerfMeasurer/requests"
)

func main() {
	var jsonfile = flag.String("jsonfile", "", "Name and path of the json file, takes precedence over other file flags")
	var filename = flag.String("filename", "", "Name of the file")
	var yamlfile = flag.String("yamlfile", "", "Name and path of the yaml file, takes precedence over filename flag")
	var isJson *bool = flag.Bool("json", false, "Whether the contents of the file are json")
	var isYaml *bool = flag.Bool("yaml", false, "Whether the contents of file are yaml")
	var outputFilename = flag.String("output", "./result.json", "Filename of output, by default will be ./result.json")

	var fileToParse string = ""
	flag.Parse()
	fmt.Println("tail:", flag.Args())
	fmt.Println(*yamlfile)
	if *jsonfile != "" {
		fileToParse = *jsonfile
		(*isJson) = true
	} else if *yamlfile != "" {
		fileToParse = *yamlfile
		(*isYaml) = true
	} else {
		fileToParse = *filename
	}

	if fileToParse == "" {
		log.Fatal("No filename given")
	}

	var requestFile parser.RequestsFile = parser.RequestsFile{}
	var err error = nil
	if *isJson {
		log.Println("Parsing json file")
		requestFile, err = json_handler.Parse(fileToParse)

	} else if *isYaml {
		log.Println("Parsing yaml file")

		requestFile, err = yaml_handler.Parse(fileToParse)

	} else {
		log.Fatalf("Unable to parse unknown filetype %s\n", fileToParse)
	}

	if err != nil {
		log.Fatal(err)
	}
	var output = requests.HitEndpoints(requestFile)
	err = generation.GenerateOutputFile(output, *outputFilename)

	if err != nil {
		log.Fatal(err)
	}
}
