package main

import (
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	yaml "gopkg.in/yaml.v2"
)

const usageXlsx2Yaml = `api_tool xlsx2yaml
	API定義の変換 xlsx -> yaml
Usage:
  api_tool xlsx2yaml INPUTS...
  api_tool xlsx2yaml -h | --help

Args:
	INPUTS...              入力ファイルパス（xlsx）

Options:
	-h --help                     Show this screen.
`

type Xlsx2YamlArg struct {
	Inputs []string
}

func NewXlsx2YamlArg(arguments map[string]interface{}) Xlsx2YamlArg {
	res := Xlsx2YamlArg{
		Inputs: sl(arguments["INPUTS"]),
	}
	return res
}

func RunXlsx2Yaml() {
	arguments, err := docopt.Parse(usageXlsx2Yaml, nil, true, "", false)
	if err != nil {
		panic(err)
	}

	arg := NewXlsx2YamlArg(arguments)

	for _, path := range arg.Inputs {
		_, _, _, groups := loadExcels([]string{path})
		outputYaml(strings.Replace(path, ".xlsx", ".yaml", -1), groups)
	}
}

func outputYaml(path string, yamlData interface{}) {
	outputData, err := yaml.Marshal(yamlData)
	e(err)

	file, err := os.Create(path)
	e(err)
	defer file.Close()
	file.Write(outputData)
}
