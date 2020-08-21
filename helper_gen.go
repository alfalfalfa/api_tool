package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alfalfalfa/go-jsschema"
	"github.com/flosch/pongo2"
	"github.com/jinzhu/inflection"
	"strings"
)

func (this *Property) HasAnotherType() bool {
	//fmt.Println("HasAnotherType")
	//dump(this.Type)
	return this.Type.HasAnotherType()
}

func (this *Action) NameAsSnakeCase() string {
	return CamelToSnake(this.Name)
}
func (this *Type) NameAsSnakeCase() string {
	return CamelToSnake(this.Name)
}
func (this PropertyType) NameAsSnakeCase() string {
	return CamelToSnake(this.Name)
}
func (this *Enum) NameAsSnakeCase() string {
	return CamelToSnake(this.Name)
}
func (this *Property) DummyData() interface{} {
	return this.Type.DummyData()
}

func genDummyData(properties []*Property) map[string]interface{} {
	data := make(map[string]interface{})
	for _, p := range properties {
		data[p.Name] = p.DummyData()
	}
	return data
}

func genDummyDataWithIndex(properties []*Property, i int) map[string]interface{} {
	data := make(map[string]interface{})
	for _, p := range properties {
		data[p.Name] = p.Type.DummyDataWithIndex(i)
	}
	return data
}

func (this *Enum) DummyData() int {
	return this.Members[0].Num
}

func (this *Type) DummyData() map[string]interface{} {
	return genDummyData(this.Properties)
}

func (this *Action) DummyRequestData() map[string]interface{} {
	return genDummyData(this.RequestProperties)
}

func (this *Action) DummyRequestDataAsJson() string {
	b, err := json.Marshal(this.DummyRequestData())
	e(err)
	return string(b)
}

func (this PropertyType) DummyData() interface{} {
	//fmt.Println("DummyData:" + this.ToJsonSchemaType().String())
	//dump(this)
	switch this.ToJsonSchemaType() {
	case schema.NullType:
		return nil
	case schema.IntegerType:
		if ee, ok := enumMap[this.Name]; ok {
			return ee.DummyData()
		} else {
			return 1
		}
	case schema.StringType:
		return "hoge"
	case schema.ArrayType:
		//fmt.Println("schema.ArrayType")
		res := make([]interface{}, 0)
		for i := 0; i < 3; i++ {
			res = append(res, this.GetArrayItemType().DummyDataWithIndex(i))
		}
		return res
	case schema.BooleanType:
		return true
	case schema.NumberType:
		return 1.1
	case schema.ObjectType:
		//fmt.Println("schema.ObjectType")
		return typeMap[this.Name].DummyData()
	}
	panic(fmt.Sprint("invalid type:", this.Name))
}

func (this PropertyType) DummyDataWithIndex(i int) interface{} {
	//fmt.Println("DummyDataWithIndex:" + this.ToJsonSchemaType().String())
	//dump(this)
	switch this.ToJsonSchemaType() {
	case schema.NullType:
		return nil
	case schema.IntegerType:
		if ee, ok := enumMap[this.Name]; ok {
			return ee.DummyDataWithIndex(i)
		} else {
			return i
		}
	case schema.StringType:
		return "mage" + base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(i)))
	case schema.ArrayType:
		//fmt.Println("schema.ArrayType")
		res := make([]interface{}, 0)
		for i := 0; i < 3; i++ {
			res = append(res, this.GetArrayItemType().DummyDataWithIndex(i))
		}
		return res
	case schema.BooleanType:
		return i%2 == 1
	case schema.NumberType:
		return 1.1 * float64(i)
	case schema.ObjectType:
		return typeMap[this.Name].DummyDataWithIndex(i)
	}
	panic(fmt.Sprint("invalid type:", this.Name))
}

func (this *Enum) DummyDataWithIndex(i int) int {
	return this.Members[i%len(this.Members)].Num
}

func (this *Type) DummyDataWithIndex(i int) map[string]interface{} {
	return genDummyDataWithIndex(this.Properties, i)
}

func (this PropertyType) DummyDataAsJson() string {
	//fmt.Println("DummyDataAsJson")
	//dump(this)
	b, err := json.Marshal(this.DummyData())
	e(err)
	return string(b)
}

// filters
func init() {
	pongo2.RegisterFilter("snake", filterSnake)
	pongo2.RegisterFilter("camel_lower", filterLowerCamel)
	pongo2.RegisterFilter("camel_upper", filterUpperCamel)
	pongo2.RegisterFilter("singular", filterSingular)
	pongo2.RegisterFilter("plural", filterPlural)
	pongo2.RegisterFilter("trim_suffix", filterTrimSuffix)
	inflection.AddIrregular("bonus", "bonuses")
}
func filterSnake(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//fmt.Println("filterSnake", in.String())
	return pongo2.AsValue(CamelToSnake(in.String())), nil
}
func filterLowerCamel(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//fmt.Println("filterLowerCamel", in.String())
	return pongo2.AsValue(SnakeToLowerCamel(in.String())), nil
}
func filterUpperCamel(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//fmt.Println("filterUpperCamel", in.String())
	return pongo2.AsValue(SnakeToCamel(in.String())), nil
}
func filterSingular(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//fmt.Println("filterSingular", in.String())
	return pongo2.AsValue(inflection.Singular(in.String())), nil
}
func filterPlural(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	//fmt.Println("filterPlural", in.String())
	return pongo2.AsValue(inflection.Plural(in.String())), nil
}
func filterTrimSuffix(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(strings.TrimSuffix(in.String(), param.String())), nil
}
