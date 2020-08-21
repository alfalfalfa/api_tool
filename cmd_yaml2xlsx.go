package main

import (
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/gobuffalo/packr"
	"github.com/tealeg/xlsx"
)

const usageYaml2Xlsx = `api_tool yaml2xlsx
	API定義の変換 yaml -> xlsx
Usage:
  api_tool yaml2xlsx INPUTS...
  api_tool yaml2xlsx -h | --help

Args:
	INPUTS...              入力ファイルパス（yaml）

Options:
	-h --help                     Show this screen.
`

type Yaml2XlsxArg struct {
	Inputs []string
}

func NewYaml2XlsxArg(arguments map[string]interface{}) Yaml2XlsxArg {
	res := Yaml2XlsxArg{
		Inputs: sl(arguments["INPUTS"]),
	}
	return res
}

func RunYaml2Xlsx() {
	arguments, err := docopt.Parse(usageYaml2Xlsx, nil, true, "", false)
	if err != nil {
		panic(err)
	}

	arg := NewYaml2XlsxArg(arguments)

	for _, path := range arg.Inputs {
		_, _, _, groups := loadYamls([]string{path})
		outputXlsx(strings.Replace(path, ".yaml", ".xlsx", -1), groups)
	}
}

func outputXlsx(path string, groups Groups) {
	bytes, err := packr.NewBox("./templates").MustBytes("template.xlsx")
	e(err)
	xlFile, err := xlsx.OpenBinary(bytes)
	e(err)

	enumTemplateSheat := xlFile.Sheet["enum"]
	typeTemplateSheat := xlFile.Sheet["type"]
	actionTemplateSheat := xlFile.Sheet["action"]
	typeMapSheat := xlFile.Sheet["_typemap"]

	for _, group := range groups {
		if 0 < len(group.Enums) {
			sheet, err := xlFile.AppendSheet(*enumTemplateSheat, "enum_"+group.Name)
			e(err)
			writeEnumsToSheet(sheet, group.Enums)
		}
		if 0 < len(group.Types) {
			sheet, err := xlFile.AppendSheet(*typeTemplateSheat, "type_"+group.Name)
			e(err)
			writeTypesToSheet(sheet, group.Types)
		}
		if 0 < len(group.Actions) {
			sheet, err := xlFile.AppendSheet(*actionTemplateSheat, "action_"+group.Name)
			e(err)
			writeActionsToSheet(sheet, group.Actions)
		}
	}

	// テンプレートシート削除
	xlFile.Sheets = xlFile.Sheets[4:]

	// _typemapシート末尾に
	xlFile.Sheets = append(xlFile.Sheets, typeMapSheat)

	e(xlFile.Save(path))
}

func writeEnumsToSheet(sheet *xlsx.Sheet, enums []*Enum) {
	for _, enum := range enums {
		row := sheet.AddRow()
		row.AddCell().SetValue(enum.Description)
		row.AddCell().SetValue(enum.Name)

		for _, member := range enum.Members {
			row.AddCell().SetValue(member.Name)
			row.AddCell().SetValue(member.Ordinal)
			row.AddCell().SetValue(member.DisplayName)
			row.AddCell().SetValue(member.Description)
			if member.Comments != nil {
				for _, c := range member.Comments {
					row.AddCell().SetValue(c)
				}
			}

			row = sheet.AddRow()
			row.AddCell()
			row.AddCell()
		}
	}
}

func writeTypesToSheet(sheet *xlsx.Sheet, types []*Type) {
	for _, typee := range types {
		row := sheet.AddRow()
		row.AddCell().SetValue(typee.Description)
		row.AddCell().SetValue(typee.Name)
		for i, prop := range typee.Properties {
			row.AddCell().SetValue(prop.Name)
			row.AddCell().SetValue(prop.Type)
			row.AddCell().SetValue(prop.Description)

			if typee.Comments != nil {
				if comments, ok := typee.Comments[i]; ok {
					for _, c := range comments {
						row.AddCell().SetValue(c)
					}
				}
			}

			row = sheet.AddRow()
			row.AddCell()
			row.AddCell()
		}
	}
}

func writeActionsToSheet(sheet *xlsx.Sheet, actions []*Action) {
	for _, action := range actions {
		for i := 0; i < len(action.RequestProperties) || i < len(action.ResponseProperties); i++ {
			row := sheet.AddRow()
			if i == 0 {
				row.AddCell().SetValue(action.Description)
				row.AddCell().SetValue(action.Name)
			} else {
				row.AddCell()
				row.AddCell()
			}

			if i < len(action.RequestProperties) {
				prop := action.RequestProperties[i]
				row.AddCell().SetValue(prop.Name)
				row.AddCell().SetValue(prop.Type)
				row.AddCell().SetValue(prop.Description)
			} else {
				row.AddCell()
				row.AddCell()
				row.AddCell()
			}

			if i < len(action.ResponseProperties) {
				prop := action.ResponseProperties[i]
				row.AddCell().SetValue(prop.Name)
				row.AddCell().SetValue(prop.Type)
				row.AddCell().SetValue(prop.Description)
			}

			if action.Comments != nil {
				if comments, ok := action.Comments[i]; ok {
					for _, c := range comments {
						row.AddCell().SetValue(c)
					}
				}
			}
		}
		sheet.AddRow()
	}
}
