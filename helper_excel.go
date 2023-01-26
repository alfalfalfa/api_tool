package main

import (
	"errors"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

func checkCommentRow(row *xlsxRow) bool {
	if len(row.Cells) == 0 {
		return true
	}

	v := getCellString(row, 0)
	if strings.HasPrefix(v, "//") {
		return true
	}
	return false
}

func checkCommentStartRow(row *xlsxRow) bool {
	v := getCellString(row, 0)
	if strings.HasPrefix(v, "//==") {
		return true
	}
	return false
}

func getCellString(row *xlsxRow, i int) string {
	if len(row.Cells) <= i {
		return ""
	}
	return strings.TrimSpace(row.Cells[i].Value)
}

func getCellStrings(row *xlsxRow, i int) []string {
	res := make([]string, 0)
	lastNotEmptyIndex := 0
	for ; i < len(row.Cells); i++ {
		v := getCellString(row, i)
		res = append(res, v)
		if v != "" {
			lastNotEmptyIndex = len(res)
			lastNotEmptyIndex = len(res)
		}
	}

	if lastNotEmptyIndex == 0 {
		return nil
	}
	return res[:lastNotEmptyIndex]
}

func getCellInt(row *xlsxRow, i int) (int, error) {
	if len(row.Cells) <= i {
		return 0, errors.New("index out of range")
	}
	return row.Cells[i].Int()
}
func getCellIntEmptyZero(row *xlsxRow, i int) int {
	if len(row.Cells) <= i {
		panic("index out of range")
	}
	res, err := row.Cells[i].Int()
	if err != nil {
		return 0
	}
	return res
}

func getCellBool(row *xlsxRow, i int) bool {
	return getCellString(row, i) != ""
}

func separateModifier(v string) (string, string) {
	var i int
	i = strings.LastIndex(v, ".")
	if i != -1 {
		return v[:i+1], v[i+1:]
	}
	i = strings.LastIndex(v, "/")
	if i != -1 {
		return v[:i+1], v[i+1:]
	}
	return "", v
}

type xlsxSheet struct {
	f    *excelize.File
	Name string
	Rows []*xlsxRow
}

func newXlsxSheet(f *excelize.File, Name string) *xlsxSheet {
	s := &xlsxSheet{f, Name, make([]*xlsxRow, 0)}
	rows, err := f.GetRows(Name)
	e(err)
	for _, row := range rows {
		xr := &xlsxRow{make([]*xlsxCell, 0)}
		for _, colCell := range row {
			xc := &xlsxCell{colCell}
			xr.Cells = append(xr.Cells, xc)
		}
		s.Rows = append(s.Rows, xr)
	}

	return s
}

//	func (s *xlsxSheet) GetRows() {
//		rows, err := s.f.Rows(s.Name)
//		e(err)
//		for rows.Next() {
//
//		}
//		l := make([]*excelize.Rows)
//	}
type xlsxRow struct {
	Cells []*xlsxCell
}
type xlsxCell struct {
	Value string
}

func (c *xlsxCell) Int() (int, error) {
	return strconv.Atoi(c.Value)
}

// Excel -> 中間データ
func loadExcels(pathes []string) ([]*Enum, TypeList, []*Action, Groups) {
	groups := Groups(make([]*Group, 0))
	enumSheets := make([]*xlsxSheet, 0)
	typeSheets := make([]*xlsxSheet, 0)
	actionSheets := make([]*xlsxSheet, 0)

	for _, path := range pathes {
		xlFile, err := excelize.OpenFile(path)
		e(err)

		for _, sheetName := range xlFile.GetSheetList() {
			// _始まりのシートは無視する
			if strings.HasPrefix(sheetName, "_") {
				continue
			}

			// 並び順をブック指定順、シート定義順にしたいので事前生成
			groups.findOrCreate(GroupNameFromSheetName(sheetName))

			// シート名prefixで仕分け
			if strings.HasPrefix(sheetName, "enum_") {
				enumSheets = append(enumSheets, newXlsxSheet(xlFile, sheetName))
				continue
			}
			if strings.HasPrefix(sheetName, "type_") {
				typeSheets = append(typeSheets, newXlsxSheet(xlFile, sheetName))
				continue
			}
			if strings.HasPrefix(sheetName, "action_") {
				actionSheets = append(actionSheets, newXlsxSheet(xlFile, sheetName))
				continue
			}
		}
	}

	enums := loadEnums(enumSheets)
	types := loadTypes(typeSheets)
	actions := loadActions(actionSheets)

	// groupで仕分け
	for _, enum := range enums {
		groups.AddEnum(enum)
	}
	for _, typee := range types {
		groups.AddType(typee)
	}
	for _, action := range actions {
		groups.AddAction(action)
	}

	return enums, types, actions, groups
}

// 全enum定義を複数シートから読み込み
func loadEnums(enumSheets []*xlsxSheet) []*Enum {
	var err error
	enums := make([]*Enum, 0)
	// enum定義ごとに行をまとめる
	rowGroup := make(map[*Enum][]*xlsxRow)
	for _, sheet := range enumSheets {
		group := GroupNameFromSheetName(sheet.Name)
		var currentEnum *Enum
		var currentRows []*xlsxRow
		for ri, r := range sheet.Rows {
			// skip header
			if ri == 0 {
				continue
			}
			if checkCommentStartRow(r) {
				break
			}
			if checkCommentRow(r) {
				continue
			}

			// get enum values
			nameTmp := getCellString(r, 1)
			if nameTmp != "" {
				if currentEnum != nil {
					rowGroup[currentEnum] = currentRows
					enums = append(enums, currentEnum)
				}
				currentEnum = NewEnum()
				currentEnum.Group = group
				currentEnum.Modifier, currentEnum.Name = separateModifier(nameTmp)
				currentEnum.Description = getCellString(r, 0)
				currentRows = make([]*xlsxRow, 0)
				currentRows = append(currentRows, r)
			} else {
				// correct property row
				if getCellString(r, 2) != "" {
					currentRows = append(currentRows, r)
				}
			}
		}
		if currentEnum != nil {
			rowGroup[currentEnum] = currentRows
			enums = append(enums, currentEnum)
		}
	}

	// enumごとのメンバー読み込み
	for t, rows := range rowGroup {
		num := -1
		for _, r := range rows {
			member := &EnumMember{
				Name:        getCellString(r, 2),
				DisplayName: getCellString(r, 4),
				Description: getCellString(r, 5),
				Comments:    getCellStrings(r, 6),
			}
			numStr := getCellString(r, 3)
			// 番号指定されていないなら連番とする
			if numStr == "" {
				num++
				member.Ordinal = num
			} else {
				member.Ordinal, err = getCellInt(r, 3)
				e(err)
				num = member.Ordinal
			}

			t.Members = append(t.Members, member)
		}
	}

	return enums
}

// 全型定義を複数シートから読み込み
func loadTypes(typeSheets []*xlsxSheet) TypeList {
	types := TypeList(make([]*Type, 0))
	// 型定義ごとに行をまとめる
	rowGroup := make(map[*Type][]*xlsxRow)
	for _, sheet := range typeSheets {
		group := GroupNameFromSheetName(sheet.Name)
		var currentType *Type
		var currentRows []*xlsxRow
		for ri, r := range sheet.Rows {
			// skip header
			if ri == 0 {
				continue
			}
			if checkCommentStartRow(r) {
				break
			}
			if checkCommentRow(r) {
				continue
			}

			// get type values
			nameTmp := getCellString(r, 1)
			if nameTmp != "" {
				if currentType != nil {
					rowGroup[currentType] = currentRows
					types = append(types, currentType)
				}
				currentType = NewType()
				currentType.Group = group
				currentType.Modifier, currentType.Name = separateModifier(nameTmp)
				currentType.Description = getCellString(r, 0)
				currentRows = make([]*xlsxRow, 0)
				currentRows = append(currentRows, r)
			} else {
				// correct property row
				if getCellString(r, 2) != "" {
					currentRows = append(currentRows, r)
				}
			}
		}
		if currentType != nil {
			rowGroup[currentType] = currentRows
			types = append(types, currentType)
		}
	}

	// 型ごとのプロパティ読み込み
	for t, rows := range rowGroup {
		t.Comments = make(map[int][]string)
		for ri, r := range rows {
			t.Properties = append(t.Properties, &Property{
				Name:        getCellString(r, 2),
				Type:        PropertyType(getCellString(r, 3)),
				Description: getCellString(r, 4),
			})

			comments := getCellStrings(r, 5)
			if comments != nil {
				t.Comments[ri] = comments
			}
		}
	}

	return types
}

// 全アクション定義を複数シートから読み込み
func loadActions(actionSheets []*xlsxSheet) []*Action {
	actions := make([]*Action, 0)
	// アクション定義ごとに行をまとめる
	rowGroup := make(map[*Action][]*xlsxRow)
	for _, sheet := range actionSheets {
		group := GroupNameFromSheetName(sheet.Name)
		var currentAction *Action
		var currentRows []*xlsxRow
		for ri, r := range sheet.Rows {
			// skip header
			if ri == 0 {
				continue
			}
			if checkCommentStartRow(r) {
				break
			}
			if checkCommentRow(r) {
				continue
			}

			// get action values
			nameTmp := getCellString(r, 1)
			if nameTmp != "" {
				if currentAction != nil {
					rowGroup[currentAction] = currentRows
					actions = append(actions, currentAction)
				}
				currentAction = NewAction()
				currentAction.Group = group
				currentAction.Name = nameTmp
				currentAction.Description = getCellString(r, 0)
				currentRows = make([]*xlsxRow, 0)
				currentRows = append(currentRows, r)
			} else {
				// correct property row
				if getCellString(r, 2) != "" || getCellString(r, 5) != "" {
					currentRows = append(currentRows, r)
				}
			}
		}
		if currentAction != nil {
			rowGroup[currentAction] = currentRows
			actions = append(actions, currentAction)
		}
	}

	// アクションごとのプロパティ読み込み
	for a, rows := range rowGroup {
		a.Comments = make(map[int][]string)
		for ri, r := range rows {
			requestPropertyName := getCellString(r, 2)
			if requestPropertyName != "" {
				a.RequestProperties = append(a.RequestProperties, &Property{
					Name:        requestPropertyName,
					Type:        PropertyType(getCellString(r, 3)),
					Description: getCellString(r, 4),
				})
			}
			responsePropertyName := getCellString(r, 5)
			if responsePropertyName != "" {
				a.ResponseProperties = append(a.ResponseProperties, &Property{
					Name:        responsePropertyName,
					Type:        PropertyType(getCellString(r, 6)),
					Description: getCellString(r, 7),
				})
			}
			comments := getCellStrings(r, 8)
			if comments != nil {
				a.Comments[ri] = comments
			}
		}
	}
	return actions
}
