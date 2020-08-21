package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/alfalfalfa/go-jsschema"
)

func (p *Property) ToSchema(index int) *schema.Schema {
	ps := schema.New()
	ps.Extras = make(map[string]interface{})
	ps.Description = p.Description
	ps.Type = append(ps.Type, p.Type.ToJsonSchemaType())
	ps.Format = schema.Format(p.Format)
	ps.Extras["order"] = index
	return ps
}

// FIXME 祖先までたどるか、プロパティに展開してしまうか
func setInheritance(s *schema.Schema, baseType string) {
	is := schema.New()
	is.Reference = "#/definitions/" + baseType
	s.AllOf = append(s.AllOf, is)
}

// dumpSchema=================================================
func dumpSchema(s *schema.Schema) {
	// ルート要素
	fmt.Println(s.SchemaRef)
	fmt.Println(s.Type)
	fmt.Println(s.ID)
	fmt.Println(s.Title)
	fmt.Println(s.Description)
	fmt.Println(s.Extras)

	re := regexp.MustCompile("\\{.*\\}")
	for name, pdef := range s.Properties {
		if "ordered_activity_request_sheet" != name {
			continue
		}

		ss, err := pdef.Resolve(nil)
		if err != nil {
			log.Printf("failed to resolv ref: %s", err)
			return
		}
		fmt.Println("\t", name, ss.Title, ss.Description, len(ss.Links), len(pdef.Links))
		//fmt.Println("\t\t", pdef.Extras)
		//fmt.Println("\t\t", ss.Extras)
		for _, link := range ss.Links {
			//re.ReplaceAllString(link.Href, "hoge")
			fmt.Println("\t\t", re.ReplaceAllString(link.Href, ":param"))
			fmt.Println("\t\t\t", link.Method, link.Title, ",", link.Description)
			//fmt.Println("\t\t\t",link.Extras)
			//delete(link.Extras, "identity")
			//if len(link.Extras)> 0{
			//	panic("1")
			//}
			if link.Schema != nil {
				//fmt.Println("\t\t\t", "request:", link.Schema.Title, ", ", link.Description, len(link.Schema.Properties), link.Schema.Reference)
				fmt.Println("\t\t\t", "request:")
				dumpSchemaProperties(4, link.Schema)
			}
			if link.TargetSchemaOrParent() != nil {
				//fmt.Println("\t\t\t", "request:", link.Schema.Title, ", ", link.Description, len(link.Schema.Properties), link.Schema.Reference)
				fmt.Println("\t\t\t", "response:")
				dumpSchemaProperties(4, link.TargetSchemaOrParent())
			}

			//fmt.Println("\t\t\t",link.TargetSchema)
		}

		//fmt.Println(pdef.IsResolved())
		//fmt.Println(ss.IsResolved())

		// Do what you will with `pdef`, which contain
		// Schema information for `name` property
		_ = name
		_ = pdef
	}

	////s.Resolve()
	//for k,v := range s.Definitions{
	//	fmt.Println(k, v.Description)
	//}

	//// You can also validate an arbitrary piece of data
	//var p interface{} // initialize using json.Unmarshal...
	//v := validator.New(s)
	//if err := v.Validate(p); err != nil {
	//	log.Printf("failed to validate data: %s", err)
	//}

	dump(s)
}

func dumpSchemaProperties(i int, s *schema.Schema) {
	var err error
	if !s.IsResolved() {
		fmt.Println(indent(i), dumpSchemaSummary(s), ", Root:", s.Root().Title)
		s, err = s.Resolve(nil)
		if err != nil {
			log.Printf("failed to resolv ref: %s", err)
			return
		}
		fmt.Println(indent(i), "->", dumpSchemaSummary(s))
	} else {
		fmt.Println(indent(i), dumpSchemaSummary(s))
	}

	for name, pdef := range s.Properties {
		fmt.Println(indent(i), "Property Name:", name)
		dumpSchemaProperties(i+1, pdef)
		//fmt.Println(indent(i), pdef.Reference, len(pdef.Properties))
		//ss, err := pdef.Resolve(nil)
		//fmt.Println(indent(i), name, ss.Title, ss.Description, ss.Type)
	}
}

func indent(indent int) string {
	res := ""
	for i := 0; i < indent; i++ {
		res += "\t"
	}
	return res
}

func dumpSchemaSummary(s *schema.Schema) string {
	res := "Schema "
	summary := make([]string, 0)
	if s.Title != "" {
		summary = append(summary, "Title:"+s.Title)
	}
	if len(s.Type) > 0 {
		summary = append(summary, "Type:"+fmt.Sprint(s.Type))
	}
	if s.Description != "" {
		summary = append(summary, "Description:"+s.Description)
	}
	if s.Reference != "" {
		summary = append(summary, "Reference:"+s.Reference)
	}

	return res + strings.Join(summary, ", ")
}
