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

const usageGenSingle = `api_tool gen-single
	API定義とpongo2テンプレートからテキスト生成
Usage:
  api_tool gen-single [--only=<OUTPUT_GROUPS>] [--overwrite=<OVERWRITE_MODE>] [--arg=<ARGUMENTS>] <OUTPUT_PATH_PETTERN> <TEMPLATE_PATH> INPUTS...
  api_tool gen-single -h | --help

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
	--only=<OUTPUT_GROUPS>        出力するグループ名をカンマ区切りで複数指定
	-h --help                     Show this screen.
`

type GenSingleArg struct {
	Inputs            []string
	OutputGroups      []string
	OutputPathPetterm string
	TemplatePath      string
	OverwriteMode     string
	Arguments         map[string]string
}

func (this GenSingleArg) IsClear() bool {
	return this.OverwriteMode == "clear"
}
func (this GenSingleArg) IsForce() bool {
	return this.OverwriteMode == "force"
}
func (this GenSingleArg) IsSkip() bool {
	return this.OverwriteMode == "skip"
}

func NewGenSingleArg(arguments map[string]interface{}) GenSingleArg {

	res := GenSingleArg{
		Inputs:            sl(arguments["INPUTS"]),
		OutputGroups:      slc(arguments["--only"]),
		OutputPathPetterm: s(arguments["<OUTPUT_PATH_PETTERN>"]),
		TemplatePath:      s(arguments["<TEMPLATE_PATH>"]),
		OverwriteMode:     s(arguments["--overwrite"]),
		Arguments:         make(map[string]string),
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

func RunGenSingle() {
	arguments, err := docopt.Parse(usageGenSingle, nil, true, "", false)
	if err != nil {
		panic(err)
	}

	arg := NewGenSingleArg(arguments)

	//fmt.Println(arg)

	enums, types, actions, groups := load(arg.Inputs, arg.OutputGroups)

	//dump(types)
	//dump(actions)

	pongo2.DefaultSet.Options.TrimBlocks = true
	pongo2.DefaultSet.Options.LStripBlocks = true
	context := pongo2.Context{
		"enums":   enums,
		"types":   types,
		"actions": actions,
		"groups":  groups,
	}

	// 追加の引数を設定
	for k, v := range arg.Arguments {
		context[k] = v
	}

	outputPathTpl, err := pongo2.DefaultSet.FromString(arg.OutputPathPetterm)
	e(err)
	outputPath, err := outputPathTpl.Execute(context)
	e(err)

	if arg.IsClear() {
		os.RemoveAll(path.Dir(outputPath))
	}

	os.MkdirAll(path.Dir(outputPath), os.ModePerm)
	tpl, err := pongo2.DefaultSet.FromFile(arg.TemplatePath)
	e(err)
	//tpl.Options.TrimBlocks = true
	//tpl.Options.LStripBlocks = true
	res, err := tpl.Execute(context)
	e(err)
	//tpl.ExecuteWriter(pongo2.Context{"page": page}, w)

	// 先頭・末尾の改行を削除し、最終行は殻行を作る
	res = strings.TrimSpace(res) + "\n"

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
