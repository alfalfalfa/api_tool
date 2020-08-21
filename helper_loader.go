package main

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var typeMap map[string]*Type
var enumMap map[string]*Enum

func load(pathes []string, outputGroups []string) ([]*Enum, TypeList, []*Action, []*Group) {
	enums, types, actions, groups := loadAny(pathes)

	// 関連解消用index
	typeMap = make(map[string]*Type)
	for _, t := range types {
		typeMap[t.Name] = t
	}
	enumMap = make(map[string]*Enum)
	for _, e := range enums {
		enumMap[e.Name] = e
	}

	//fmt.Println(outputGroups)
	// outputGroupsでフィルタリング
	if outputGroups != nil && len(outputGroups) > 0 {
		filtered_enums := make([]*Enum, 0)
		filtered_types := TypeList(make([]*Type, 0))
		filtered_actions := make([]*Action, 0)

		for _, enum := range enums {
			for _, group := range outputGroups {
				if group == enum.Group {
					filtered_enums = append(filtered_enums, enum)
				}
			}
		}
		for _, typee := range types {
			for _, group := range outputGroups {
				if group == typee.Group {
					filtered_types = append(filtered_types, typee)
				}
			}
		}
		for _, action := range actions {
			for _, group := range outputGroups {
				if group == action.Group {
					filtered_actions = append(filtered_actions, action)
				}
			}
		}

		// TODO groupのフィルタリング必要？

		return filtered_enums, filtered_types, filtered_actions, groups
	}

	return enums, types, actions, groups
}

// YAML|Xlsx -> 中間データ
func loadAny(pathes []string) ([]*Enum, TypeList, []*Action, Groups) {
	if strings.HasSuffix(pathes[0], "xlsx") {
		return loadExcels(pathes)
	}
	if strings.HasSuffix(pathes[0], "yaml") {
		return loadYamls(pathes)
	}
	panic("invalid filetype:" + pathes[0])
}

// YAML -> 中間データ
func loadYamls(pathes []string) ([]*Enum, TypeList, []*Action, Groups) {
	enums := make([]*Enum, 0)
	types := TypeList(make([]*Type, 0))
	actions := make([]*Action, 0)
	groups := Groups(make([]*Group, 0))
	//groupsList = make(Groups[], len(pathes))
	for _, path := range pathes {
		buf, err := ioutil.ReadFile(path)
		e(err)
		var tmpGroups Groups
		err = yaml.Unmarshal(buf, &tmpGroups)
		e(err)

		for _, group := range tmpGroups {
			// 並び順をブック指定順、シート定義順にしたいので事前生成
			groups.findOrCreate(group.Name)
			enums = append(enums, group.Enums...)
			types = append(types, group.Types...)
			actions = append(actions, group.Actions...)
		}
	}

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
