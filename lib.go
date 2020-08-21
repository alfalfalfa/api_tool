package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func e(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func s(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}
func sl(v interface{}) []string {
	if v == nil {
		return make([]string, 0)
	}
	if sl, ok := v.([]string); ok {
		return sl
	}
	return make([]string, 0)
}

// カンマ区切り
func slc(v interface{}) []string {
	res := make([]string, 0)
	if v == nil {
		return res
	}
	tmp1 := fmt.Sprint(v)
	tmp2 := strings.Split(tmp1, ",")
	for _, str := range tmp2 {
		str = strings.TrimSpace(str)
		if str != "" {
			res = append(res, str)
		}
	}
	return res
}
func b(v interface{}) bool {
	if v == nil {
		return false
	}
	if bb, ok := v.(bool); ok {
		return bb
	}
	return false
}

func i(v interface{}) int {
	if v == nil {
		return 0
	}
	if ii, ok := v.(int); ok {
		return ii
	}
	if sv, ok := v.(string); ok {
		ii, err := strconv.Atoi(sv)
		if err != nil {
			return 0
		}
		return ii
	}
	return 0
}

func dump(v interface{}) {
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
	//fmt.Sprintln(string(j))
}

func dump2str(v interface{}) string {
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(j))
	return string(j)
}

func dump2yaml(v interface{}) {
	j, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
	//fmt.Sprintln(string(j))
}
