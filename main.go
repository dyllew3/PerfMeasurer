package main

import (
	"fmt"

	yaml_handler "github.com/dyllew3/PerfMeasurer/parser/yaml"
)

func main() {
	fmt.Println("Got here")
	result, err := yaml_handler.Parse("./example.yml")
	fmt.Println(err)
	fmt.Println(result)
}
