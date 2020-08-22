package main

import (
	"fmt"

	"os"
	"path"

	"strings"

	"io/ioutil"

	"github.com/docopt/docopt-go"
	"github.com/flosch/pongo2"
)

const usageGenMultiple = `api_tool gen-multiple
	API定義とpongo2テンプレートから複数のテキスト生成
Usage:
  api_tool gen-multiple (--target=<TARGET>) [--only=<OUTPUT_GROUPS>] [--overwrite=<OVERWRITE_MODE>] [--arg=<ARGUMENTS>] <OUTPUT_PATH_PETTERN> <TEMPLATE_PATH> INPUTS...
  api_tool gen-multiple -h | --help

Args:
	<OUTPUT_PATH_PETTERN>  (必須)出力ファイルパスパターン pongo2
	<TEMPLATE_PATH>        (必須)テンプレートファイルパス
	INPUTS...              入力ファイルパス（xlsx）

Options:
	--overwrite=<OVERWRITE_MODE>  Overwrite behavior [default:force]
							        force 上書き
									skip  存在していたらスキップ
									clear 出力ディレクトリを削除
	--arg=<ARGUMENTS>             Additional arguments to pongo2
							        key1:value1,key2:value2
	--target=<TARGET>             multiple target
							        action 各action定義からそれぞれファイル出力
							        type 各type定義からそれぞれファイル出力
							        enum 各enum定義からそれぞれファイル出力
	--only=<OUTPUT_GROUPS>        出力するグループ名をカンマ区切りで複数指定
	-h --help                     Show this screen.
`

type GenMultipleArg struct {
	Inputs            []string
	OutputGroups      []string
	OutputPathPetterm string
	TemplatePath      string
	OverwriteMode     string
	Arguments         map[string]string
	Target            string
}

func (this GenMultipleArg) IsClear() bool {
	return this.OverwriteMode == "clear"
}
func (this GenMultipleArg) IsForce() bool {
	return this.OverwriteMode == "force"
}
func (this GenMultipleArg) IsSkip() bool {
	return this.OverwriteMode == "skip"
}

func NewGenMultipleArg(arguments map[string]interface{}) GenMultipleArg {
	res := GenMultipleArg{
		Inputs:            sl(arguments["INPUTS"]),
		OutputGroups:      slc(arguments["--only"]),
		OutputPathPetterm: s(arguments["<OUTPUT_PATH_PETTERN>"]),
		TemplatePath:      s(arguments["<TEMPLATE_PATH>"]),
		OverwriteMode:     s(arguments["--overwrite"]),
		Arguments:         make(map[string]string),
		Target:            s(arguments["--target"]),
	}
	arg := s(arguments["--arg"])
	if arg != "" {
		for _, v := range strings.Split(arg, ",") {
			tmp := strings.Split(v, ":")
			res.Arguments[tmp[0]] = tmp[1]
		}
	}

	return res
}

func RunGenMultiple() {
	arguments, err := docopt.Parse(usageGenMultiple, nil, true, "", false)
	if err != nil {
		panic(err)
	}

	arg := NewGenMultipleArg(arguments)

	enums, types, actions, _ := load(arg.Inputs, arg.OutputGroups)

	//dump(types)
	//dump(actions)
	pongo2.DefaultSet.Options.TrimBlocks = true
	pongo2.DefaultSet.Options.LStripBlocks = true
	outputPathTpl, err := pongo2.DefaultSet.FromString(arg.OutputPathPetterm)
	e(err)
	tpl, err := pongo2.DefaultSet.FromFile(arg.TemplatePath)
	e(err)

	outputs := make(map[string]string)

	switch arg.Target {
	case "action":
		for _, action := range actions {
			context := pongo2.Context{
				"enums":   enums,
				"types":   types,
				"actions": actions,
				"action":  action,
			}
			// 追加の引数を設定
			for k, v := range arg.Arguments {
				context[k] = v
			}
			outputPath, res := renderTemplate(context, outputPathTpl, tpl)
			outputs[outputPath] = res
		}
	case "type":
		for _, t := range types {
			context := pongo2.Context{
				"enums":   enums,
				"types":   types,
				"actions": actions,
				"type":    t,
			}
			// 追加の引数を設定
			for k, v := range arg.Arguments {
				context[k] = v
			}
			outputPath, res := renderTemplate(context, outputPathTpl, tpl)
			outputs[outputPath] = res
		}
	case "enum":
		for _, e := range enums {
			context := pongo2.Context{
				"enums":   enums,
				"types":   types,
				"actions": actions,
				"enum":    e,
			}
			// 追加の引数を設定
			for k, v := range arg.Arguments {
				context[k] = v
			}
			outputPath, res := renderTemplate(context, outputPathTpl, tpl)
			outputs[outputPath] = res
		}
	}

	dirs := make(map[string]bool)
	for outputPath, _ := range outputs {
		dirs[path.Dir(outputPath)] = true
	}

	for outputDir, _ := range dirs {
		if arg.IsClear() {
			os.RemoveAll(outputDir)
		}
		os.MkdirAll(outputDir, os.ModePerm)
	}

	for outputPath, res := range outputs {
		if arg.IsSkip() {
			_, err := os.Stat(outputPath)
			if err != nil {
				ioutil.WriteFile(outputPath, []byte(res), os.ModePerm)
				fmt.Println("write:", outputPath)
			} else {
				fmt.Println("skip:", outputPath)
			}
		} else {
			ioutil.WriteFile(outputPath, []byte(res), os.ModePerm)
			fmt.Println("write:", outputPath)
		}
	}
}

func renderTemplate(context pongo2.Context, outputPathTpl, tpl *pongo2.Template) (string, string) {
	outputPath, err := outputPathTpl.Execute(context)
	e(err)

	res, err := tpl.Execute(context)
	e(err)

	// 先頭・末尾の改行を削除し、最終行は殻行を作る
	res = strings.TrimSpace(res) + "\n"

	return outputPath, res
}
