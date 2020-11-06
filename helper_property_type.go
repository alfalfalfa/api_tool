package main

import (
	"errors"

	"strings"

	"github.com/alfalfalfa/go-jsschema"
)

func (this PropertyType) ToJsonSchemaType() (t schema.PrimitiveType) {
	if this.IsArray() {
		t = schema.ArrayType
		return
	}
	switch string(this) {
	case "null":
		t = schema.NullType
	case "nil":
		t = schema.NullType
	case "integer":
		t = schema.IntegerType
	case "int":
		t = schema.IntegerType
	case "int8":
		t = schema.IntegerType
	case "int16":
		t = schema.IntegerType
	case "int32":
		t = schema.IntegerType
	case "int64":
		t = schema.IntegerType
	case "long":
		t = schema.IntegerType
	case "uint8":
		t = schema.IntegerType
	case "uint16":
		t = schema.IntegerType
	case "uint32":
		t = schema.IntegerType
	case "uint64":
		t = schema.IntegerType
	case "sint32":
		t = schema.IntegerType
	case "sint64":
		t = schema.IntegerType
	case "text":
		t = schema.StringType
	case "string":
		t = schema.StringType
	case "binary":
		// JsonはBase64を想定
		t = schema.StringType
	case "boolean":
		t = schema.BooleanType
	case "bool":
		t = schema.BooleanType
	case "number":
		t = schema.NumberType
	case "float32":
		t = schema.NumberType
	case "float":
		t = schema.NumberType
	case "float64":
		t = schema.NumberType
	case "double":
		t = schema.NumberType
	case "timestamp":
		t = schema.IntegerType
	default:
		if _, ok := typeMap[string(this)]; ok {
			t = schema.ObjectType
		} else if _, ok := enumMap[string(this)]; ok {
			t = schema.IntegerType
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}

func (this PropertyType) ToGoType(packageName string) (t string) {
	if this.IsArray() {
		t = "[]" + this.GetArrayItemType().ToGoType(packageName)
		return
	}
	switch string(this) {
	case "null":
		t = "nil"
	case "nil":
		t = "nil"
	case "integer":
		t = "int"
	case "int":
		t = "int"
	case "int8":
		t = "int8"
	case "int16":
		t = "int16"
	case "int32":
		t = "int32"
	case "int64":
		t = "int64"
	case "long":
		t = "int64"
	case "uint8":
		t = "uint8"
	case "uint16":
		t = "uint16"
	case "uint32":
		t = "uint32"
	case "uint64":
		t = "uint64"
	case "sint32":
		t = "int32"
	case "sint64":
		t = "int64"
	case "text":
		t = "string"
	case "string":
		t = "string"
	case "binary":
		t = "[]byte"
	case "boolean":
		t = "bool"
	case "bool":
		t = "bool"
	case "number":
		t = "float64"
	case "float32":
		t = "float32"
	case "float":
		t = "float32"
	case "float64":
		t = "float64"
	case "double":
		t = "float64"
	case "timestamp":
		t = "uint32"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			t = packageName + "." + tt.Name
		} else if ee, ok := enumMap[string(this)]; ok {
			t = packageName + "." + ee.Name
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}

func (this PropertyType) ToRubyType(moduleName string) (t string) {
	if this.IsArray() {
		t = "Array<" + this.GetArrayItemType().ToRubyType(moduleName) + ">"
		return
	}
	switch string(this) {
	case "null":
		t = "nil"
	case "nil":
		t = "nil"
	case "integer":
		t = "Integer"
	case "int":
		t = "Integer"
	case "int8":
		t = "Integer"
	case "int16":
		t = "Integer"
	case "int32":
		t = "Integer"
	case "int64":
		t = "Integer"
	case "long":
		t = "Integer"
	case "uint8":
		t = "Integer"
	case "uint16":
		t = "Integer"
	case "uint32":
		t = "Integer"
	case "uint64":
		t = "Integer"
	case "sint32":
		t = "Integer"
	case "sint64":
		t = "Integer"
	case "text":
		t = "String"
	case "string":
		t = "String"
	case "binary":
		t = "Array<Integer>"
	case "boolean":
		t = "true, false"
	case "bool":
		t = "true, false"
	case "number":
		t = "Numeric"
	case "float32":
		t = "Float"
	case "float":
		t = "Float"
	case "float64":
		t = "Float"
	case "double":
		t = "Float"
	case "timestamp":
		t = "Integer"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			t = moduleName + "::" + tt.Name
		} else if ee, ok := enumMap[string(this)]; ok {
			//TODO rubyのenum
			t = moduleName + "::" + ee.Name
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}

func (this PropertyType) ToCSType(namespace string) (t string) {
	if this.IsArray() {
		t = this.GetArrayItemType().ToCSType(namespace) + "[]"
		return
	}
	switch string(this) {
	case "null":
		t = "null"
	case "nil":
		t = "null"
	case "integer":
		t = "int"
	case "int":
		t = "int"
	case "int8":
		t = "sbyte"
	case "int16":
		t = "short"
	case "int32":
		t = "int"
	case "int64":
		t = "long"
	case "long":
		t = "long"
	case "uint8":
		t = "byte"
	case "uint16":
		t = "ushort"
	case "uint32":
		t = "uint"
	case "uint64":
		t = "ulong"
	case "sint32":
		t = "int"
	case "sint64":
		t = "long"
	case "text":
		t = "string"
	case "string":
		t = "string"
	case "binary":
		t = "byte[]"
	case "boolean":
		t = "bool"
	case "bool":
		t = "bool"
	case "number":
		t = "double"
	case "float32":
		t = "float"
	case "float":
		t = "float"
	case "float64":
		t = "double"
	case "double":
		t = "double"
	case "timestamp":
		t = "uint"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			t = namespace + "." + tt.Name
		} else if ee, ok := enumMap[string(this)]; ok {
			t = namespace + "." + ee.Name
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}
func (this PropertyType) ToFlowType(moduleName string) (t string) {
	if this.IsArray() {
		t = "Array<" + this.GetArrayItemType().ToFlowType(moduleName) + ">"
		return
	}
	switch string(this) {
	case "null":
		t = "null"
	case "nil":
		t = "null"
	case "integer":
		t = "number"
	case "int":
		t = "number"
	case "int8":
		t = "number"
	case "int16":
		t = "number"
	case "int32":
		t = "number"
	case "int64":
		t = "number"
	case "long":
		t = "number"
	case "uint8":
		t = "number"
	case "uint16":
		t = "number"
	case "uint32":
		t = "number"
	case "uint64":
		t = "number"
	case "sint32":
		t = "number"
	case "sint64":
		t = "number"
	case "text":
		t = "string"
	case "string":
		t = "string"
	case "binary":
		t = "ArrayBuffer"
	case "boolean":
		t = "boolean"
	case "bool":
		t = "boolean"
	case "number":
		t = "number"
	case "float32":
		t = "number"
	case "float":
		t = "number"
	case "float64":
		t = "number"
	case "double":
		t = "number"
	case "timestamp":
		t = "number"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			if moduleName == "" {
				t = tt.Name
			} else {
				t = strings.TrimSuffix(moduleName, ".") + "." + tt.Name
			}
		} else if ee, ok := enumMap[string(this)]; ok {
			if moduleName == "" {
				t = ee.Name
			} else {
				t = strings.TrimSuffix(moduleName, ".") + "." + ee.Name
			}
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}
func (this PropertyType) ToJavaType(namespace string) (t string) {
	if this.IsArray() {
		t = "List<" + this.GetArrayItemType().ToJavaClassType(namespace) + ">"
		return
	}
	switch string(this) {
	case "null":
		t = "null"
	case "nil":
		t = "null"
	case "integer":
		t = "int"
	case "int":
		t = "int"
	case "int8":
		t = "byte"
	case "int16":
		t = "short"
	case "int32":
		t = "int"
	case "int64":
		t = "long"
	case "long":
		t = "long"
	case "uint8":
		t = "short"
	case "uint16":
		t = "int"
	case "uint32":
		t = "long"
	case "uint64":
		t = "long"
	case "sint32":
		t = "int"
	case "sint64":
		t = "long"
	case "text":
		t = "String"
	case "string":
		t = "String"
	case "binary":
		t = "byte[]"
	case "boolean":
		t = "boolean"
	case "bool":
		t = "boolean"
	case "number":
		t = "double"
	case "float32":
		t = "float"
	case "float":
		t = "float"
	case "float64":
		t = "double"
	case "double":
		t = "double"
	case "timestamp":
		t = "long"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			t = namespace + tt.Name
		} else if ee, ok := enumMap[string(this)]; ok {
			t = namespace + ee.Name
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}
func (this PropertyType) ToJavaClassType(namespace string) (t string) {
	if this.IsArray() {
		t = "List<" + this.GetArrayItemType().ToJavaClassType(namespace) + ">"
		return
	}
	switch string(this) {
	case "null":
		t = "null"
	case "nil":
		t = "null"
	case "integer":
		t = "Integer"
	case "int":
		t = "Integer"
	case "int8":
		t = "Byte"
	case "int16":
		t = "Short"
	case "int32":
		t = "Integer"
	case "int64":
		t = "Long"
	case "long":
		t = "Long"
	case "uint8":
		t = "Short"
	case "uint16":
		t = "Integer"
	case "uint32":
		t = "Long"
	case "uint64":
		t = "Long"
	case "sint32":
		t = "Integer"
	case "sint64":
		t = "Long"
	case "text":
		t = "String"
	case "string":
		t = "String"
	case "binary":
		t = "Byte[]"
	case "boolean":
		t = "Boolean"
	case "bool":
		t = "Boolean"
	case "number":
		t = "Double"
	case "float32":
		t = "Float"
	case "float":
		t = "Float"
	case "float64":
		t = "Double"
	case "double":
		t = "Double"
	case "timestamp":
		t = "Long"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			t = namespace + tt.Name
		} else if ee, ok := enumMap[string(this)]; ok {
			t = namespace + ee.Name
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}
func (this PropertyType) ToTSType(namespace string) (t string) {
	if this.IsArray() {
		t = this.GetArrayItemType().ToTSType(namespace) + "[]"
		return
	}
	switch string(this) {
	case "null":
		t = "null"
	case "nil":
		t = "null"
	case "integer":
		t = "number"
	case "int":
		t = "number"
	case "int8":
		t = "number"
	case "int16":
		t = "number"
	case "int32":
		t = "number"
	case "int64":
		t = "number"
	case "long":
		t = "number"
	case "uint8":
		t = "number"
	case "uint16":
		t = "number"
	case "uint32":
		t = "number"
	case "uint64":
		t = "number"
	case "sint32":
		t = "number"
	case "sint64":
		t = "number"
	case "text":
		t = "string"
	case "string":
		t = "string"
	case "binary":
		t = "string"
	case "boolean":
		t = "boolean"
	case "bool":
		t = "boolean"
	case "number":
		t = "number"
	case "float32":
		t = "number"
	case "float":
		t = "number"
	case "float64":
		t = "number"
	case "double":
		t = "number"
	case "timestamp":
		t = "number"
	default:
		if tt, ok := typeMap[string(this)]; ok {
			t = namespace + tt.Name
		} else if ee, ok := enumMap[string(this)]; ok {
			t = namespace + ee.Name
		} else {
			e(errors.New("unknown property type: " + dump2str(this)))
		}
	}
	return
}
func (this PropertyType) IsArray() bool {
	return strings.HasPrefix(string(this), "[]") || strings.HasSuffix(string(this), "[]")
	//||	this.Name == "binary"
}

func (this PropertyType) GetArrayItemType() PropertyType {
	if !this.IsArray() {
		panic(string(this) + " is not array")
	}
	return PropertyType(strings.Replace(string(this), "[]", "", -1))
}

func (this PropertyType) HasAnotherType() bool {
	return this.ToJsonSchemaType() == schema.ObjectType || (this.ToJsonSchemaType() == schema.ArrayType && this.GetArrayItemType().HasAnotherType())
}
func (this PropertyType) GetType() *Type {
	if t, ok := typeMap[string(this)]; ok {
		return t
	} else {
		e(errors.New("unknown property type: " + dump2str(this)))
	}
	return nil
}

func (this PropertyType) IsObject() bool {
	return this.ToJsonSchemaType() == schema.ObjectType
}
func (this PropertyType) IsEnum() bool {
	_, ok := enumMap[string(this)]
	return ok
}
func (this PropertyType) GetEnum() *Enum {
	if t, ok := enumMap[string(this)]; ok {
		return t
	} else {
		e(errors.New("unknown enum type: " + dump2str(this)))
	}
	return nil
}
