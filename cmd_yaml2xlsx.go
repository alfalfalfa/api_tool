package main

import (
	"bytes"
	"strings"

	_ "embed"
	"github.com/docopt/docopt-go"
	"github.com/xuri/excelize/v2"
)

//go:embed templates/template.xlsx
var templateBytes []byte

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
	xlFile, err := excelize.OpenReader(bytes.NewReader(templateBytes))
	e(err)
	defer xlFile.Close()
	enumTemplateSheetIndex, err := xlFile.GetSheetIndex("enum")
	e(err)
	typeTemplateSheetIndex, err := xlFile.GetSheetIndex("type")
	e(err)
	actionTemplateSheetIndex, err := xlFile.GetSheetIndex("action")
	e(err)

	for _, group := range groups {
		if 0 < len(group.Enums) {
			sheet := "enum_" + group.Name
			si, err := xlFile.NewSheet(sheet)
			e(err)
			e(xlFile.CopySheet(enumTemplateSheetIndex, si))
			writeEnumsToSheet(xlFile, sheet, group.Enums)
		}
		if 0 < len(group.Types) {
			sheet := "type_" + group.Name
			si, err := xlFile.NewSheet(sheet)
			e(err)
			e(xlFile.CopySheet(typeTemplateSheetIndex, si))
			writeTypesToSheet(xlFile, sheet, group.Types)
		}
		if 0 < len(group.Actions) {
			sheet := "action_" + group.Name
			si, err := xlFile.NewSheet(sheet)
			e(err)
			e(xlFile.CopySheet(actionTemplateSheetIndex, si))
			writeActionsToSheet(xlFile, sheet, group.Actions)
		}
	}

	// _typemapシート末尾に
	e(xlFile.SetSheetName("_typemap", "__typemap"))
	typeMapSheatIndex, err := xlFile.GetSheetIndex("__typemap")
	e(err)
	newTypeMapSheatIndex, err := xlFile.NewSheet("_typemap")
	e(err)
	e(xlFile.CopySheet(typeMapSheatIndex, newTypeMapSheatIndex))

	// テンプレートシート削除
	e(xlFile.DeleteSheet("enum"))
	e(xlFile.DeleteSheet("type"))
	e(xlFile.DeleteSheet("action"))
	e(xlFile.DeleteSheet("__typemap"))

	e(xlFile.SaveAs(path))
}
func c(ri, ci int) string {
	s, err := excelize.CoordinatesToCellName(ci+1, ri+1)
	e(err)
	return s
}

type row struct {
	f     *excelize.File
	sheet string
	ri    int
	ci    int
}

func newRow(f *excelize.File, sheet string, ri int) *row {
	r := &row{
		f:     f,
		sheet: sheet,
		ri:    ri,
		ci:    -1,
	}
	return r
}
func (r *row) AddCell() *row {
	r.ci++
	return r
}
func (r *row) SetValue(v interface{}) *row {
	e(r.f.SetCellValue(r.sheet, c(r.ri, r.ci), v))
	return r
}

func writeEnumsToSheet(f *excelize.File, sheet string, enums []*Enum) {
	ri := 0
	for _, enum := range enums {
		ri++
		row := newRow(f, sheet, ri)
		row.AddCell().SetValue(enum.Description)
		row.AddCell().SetValue(enum.Modifier + enum.Name)

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

			ri++
			row = newRow(f, sheet, ri)
			row.AddCell()
			row.AddCell()
		}
	}
}
func writeTypesToSheet(f *excelize.File, sheet string, types []*Type) {
	ri := 0
	for _, typee := range types {
		ri++
		row := newRow(f, sheet, ri)
		row.AddCell().SetValue(typee.Description)
		row.AddCell().SetValue(typee.Modifier + typee.Name)
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

			ri++
			row = newRow(f, sheet, ri)
			row.AddCell()
			row.AddCell()
		}
	}
}

func writeActionsToSheet(f *excelize.File, sheet string, actions []*Action) {
	ri := 0
	for _, action := range actions {
		for i := 0; i < len(action.RequestProperties) || i < len(action.ResponseProperties); i++ {
			ri++
			row := newRow(f, sheet, ri)
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
		ri++
	}
}
