package main

import (
	"errors"
	"strings"

	"github.com/tealeg/xlsx"
)

func checkCommentRow(row *xlsx.Row) bool {
	if len(row.Cells) == 0 {
		return true
	}

	v := getCellString(row, 0)
	if strings.HasPrefix(v, "//") {
		return true
	}
	return false
}

func getCellString(row *xlsx.Row, i int) string {
	if len(row.Cells) <= i {
		return ""
	}
	return strings.TrimSpace(row.Cells[i].Value)
}

func getCellStrings(row *xlsx.Row, i int) []string {
	res := make([]string, 0)
	lastNotEmptyIndex := 0
	for ; i < len(row.Cells); i++ {
		v := getCellString(row, i)
		res = append(res, v)
		if v != "" {
			lastNotEmptyIndex = len(res)
		}
	}

	if lastNotEmptyIndex == 0 {
		return nil
	}
	return res[:lastNotEmptyIndex]
}

func getCellInt(row *xlsx.Row, i int) (int, error) {
	if len(row.Cells) <= i {
		return 0, errors.New("index out of range")
	}
	return row.Cells[i].Int()
}
func getCellIntEmptyZero(row *xlsx.Row, i int) int {
	if len(row.Cells) <= i {
		panic("index out of range")
	}
	res, err := row.Cells[i].Int()
	if err != nil {
		return 0
	}
	return res
}

func getCellBool(row *xlsx.Row, i int) bool {
	return getCellString(row, i) != ""
}

// Excel -> 中間データ
func loadExcels(pathes []string) ([]*Enum, TypeList, []*Action, Groups) {
	groups := Groups(make([]*Group, 0))
	enumSheets := make([]*xlsx.Sheet, 0)
	typeSheets := make([]*xlsx.Sheet, 0)
	actionSheets := make([]*xlsx.Sheet, 0)

	for _, path := range pathes {
		xlFile, err := xlsx.OpenFile(path)
		e(err)
		for _, sheet := range xlFile.Sheets {
			// _始まりのシートは無視する
			if strings.HasPrefix(sheet.Name, "_") {
				continue
			}

			// 並び順をブック指定順、シート定義順にしたいので事前生成
			groups.findOrCreate(GroupNameFromSheetName(sheet.Name))

			// シート名prefixで仕分け
			if strings.HasPrefix(sheet.Name, "enum_") {
				enumSheets = append(enumSheets, sheet)
				continue
			}
			if strings.HasPrefix(sheet.Name, "type_") {
				typeSheets = append(typeSheets, sheet)
				continue
			}
			if strings.HasPrefix(sheet.Name, "action_") {
				actionSheets = append(actionSheets, sheet)
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
func loadEnums(enumSheets []*xlsx.Sheet) []*Enum {
	var err error
	enums := make([]*Enum, 0)
	// enum定義ごとに行をまとめる
	rowGroup := make(map[*Enum][]*xlsx.Row)
	for _, sheet := range enumSheets {
		group := GroupNameFromSheetName(sheet.Name)
		var currentEnum *Enum
		var currentRows []*xlsx.Row
		for ri, r := range sheet.Rows {
			// skip header
			if ri == 0 {
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
				currentEnum.Name = nameTmp
				currentEnum.Title = getCellString(r, 0)
				currentRows = make([]*xlsx.Row, 0)
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
				member.Num = num
			} else {
				member.Num, err = getCellInt(r, 3)
				e(err)
				num = member.Num
			}

			t.Members = append(t.Members, member)
		}
	}

	return enums
}

// 全型定義を複数シートから読み込み
func loadTypes(typeSheets []*xlsx.Sheet) TypeList {
	types := TypeList(make([]*Type, 0))
	// 型定義ごとに行をまとめる
	rowGroup := make(map[*Type][]*xlsx.Row)
	for _, sheet := range typeSheets {
		group := GroupNameFromSheetName(sheet.Name)
		var currentType *Type
		var currentRows []*xlsx.Row
		for ri, r := range sheet.Rows {
			// skip header
			if ri == 0 {
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
				currentType.Name = nameTmp
				currentType.Title = getCellString(r, 0)
				currentType.BaseType = getCellString(r, 2)
				currentType.SheetName = sheet.Name
				currentRows = make([]*xlsx.Row, 0)
				currentRows = append(currentRows, r)
			} else {
				// correct property row
				if getCellString(r, 3) != "" {
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
				Name:        getCellString(r, 3),
				Type:        NewPropertyType(t.SheetName, t.Name, ri, 4, getCellString(r, 4)),
				Format:      getCellString(r, 5),
				Description: getCellString(r, 6),
			})

			comments := getCellStrings(r, 7)
			if comments != nil {
				t.Comments[ri] = comments
			}
		}
	}

	return types
}

// 全アクション定義を複数シートから読み込み
func loadActions(actionSheets []*xlsx.Sheet) []*Action {
	actions := make([]*Action, 0)
	// アクション定義ごとに行をまとめる
	rowGroup := make(map[*Action][]*xlsx.Row)
	for _, sheet := range actionSheets {
		group := GroupNameFromSheetName(sheet.Name)
		var currentAction *Action
		var currentRows []*xlsx.Row
		for ri, r := range sheet.Rows {
			// skip header
			if ri == 0 {
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
				currentAction.Title = getCellString(r, 0)
				currentAction.RequestBaseType = getCellString(r, 2)
				currentAction.ResponseBaseType = getCellString(r, 7)
				currentAction.SheetName = sheet.Name
				currentRows = make([]*xlsx.Row, 0)
				currentRows = append(currentRows, r)
			} else {
				// correct property row
				if getCellString(r, 3) != "" || getCellString(r, 8) != "" {
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
			requestPropertyName := getCellString(r, 3)
			if requestPropertyName != "" {
				a.RequestProperties = append(a.RequestProperties, &Property{
					Name:        requestPropertyName,
					Type:        NewPropertyType(a.SheetName, a.Name, ri, 4, getCellString(r, 4)),
					Format:      getCellString(r, 5),
					Description: getCellString(r, 6),
				})
			}
			responsePropertyName := getCellString(r, 8)
			if responsePropertyName != "" {
				a.ResponseProperties = append(a.ResponseProperties, &Property{
					Name:        responsePropertyName,
					Type:        NewPropertyType(a.SheetName, a.Name, ri, 9, getCellString(r, 9)),
					Format:      getCellString(r, 10),
					Description: getCellString(r, 11),
				})
			}
			comments := getCellStrings(r, 12)
			if comments != nil {
				a.Comments[ri] = comments
			}
		}
	}
	return actions
}
