// Copyright 2017 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// gen-accessors generates accessor methods for structs with pointer fields.
//
// It is meant to be used by go-github contributors in conjunction with the
// go generate tool before sending a PR to GitHub.
// Please see the CONTRIBUTING.md file for more information.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

func lowerCaseFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)

	first := unicode.ToLower(r[0])

	if len(r) > 1 {
		return string(append([]rune{first}, r[1:]...))
	}

	return string(first)
}

func indentTabs(tabs int, v string) string {
	pad := strings.Repeat("\t", tabs)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

var flagSrc string

func main() {
	flag.StringVar(&flagSrc, "src", "", "src file")

	flag.Parse()

	fmt.Println(flagSrc)

	tmpl := template.Must(
		template.New("source").
			Funcs(template.FuncMap{
				"lowerCaseFirst": lowerCaseFirst,
				"indentTabs":     indentTabs,
			}).
			Parse(fileTemplateSrc))

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, flagSrc, nil, 0)
	if err != nil {
		panic(err)
	}

	d := &fileTemplateData{}

	typesFound := 0

	for _, f := range f.Decls {
		dcl, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, s := range dcl.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			fmt.Println(ts.Name.Name)

			// Skip unexported identifiers.
			if !ts.Name.IsExported() {
				fmt.Printf("Struct %v is unexported; skipping.\n", ts.Name)
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			typesFound++
			if typesFound > 1 {
				panic("Found more than one exported struct in file. This is currently not supported")
			}

			d.StructTypeName = ts.Name.Name

			for _, field := range st.Fields.List {
				fieldName := field.Names[0].Name
				fmt.Println("---")
				fmt.Printf("Field name: %s\n", fieldName)
				// fmt.Printf("Field type: %s\n", reflect.TypeOf(field.Type))

				switch v := field.Type.(type) {
				case *ast.Ident:
					fmt.Println("Ident:", v.Name)

					switch v.Name {
					case "string":
						fieldData := &UnmarshalStringTemplateData{
							FieldName: fieldName,
						}
						d.Fields = append(d.Fields, fieldData)
					}

				case *ast.StarExpr:
					//fmt.Println("StarExpr:", v)
					fmt.Println("Pointer of:", v.X.(*ast.Ident).Name)

					panic("pointers not supported right now...")
					//fieldData := &UnmarshalStringTemplateData{
					//	FieldName: fieldName,
					//	FieldType: v.X.(*ast.Ident).Name,
					//}
					//d.Fields = append(d.Fields, fieldData)

				case *ast.MapType:
					//fmt.Println("MapType:", v)
					fmt.Println("Map key", v.Key.(*ast.Ident).Name)

					fieldData := &UnmarshalMapTemplateData{
						FieldName: fieldName,
					}

					switch mV := v.Value.(type) {
					case *ast.StarExpr:
						fmt.Println("Map pointer value", mV.X.(*ast.Ident).Name)
						fieldData.MapValueType = "*" + mV.X.(*ast.Ident).Name
					case *ast.Ident:
						fmt.Println("Map ident", mV.Name)
						fieldData.MapValueType = mV.Name
					default:
						panic(fmt.Sprintf("Unsupported map value type: %s", reflect.TypeOf(mV)))
					}

					d.Fields = append(d.Fields, fieldData)

				case *ast.ArrayType:

					fieldData := &UnmarshalArrayTemplateData{
						FieldName: fieldName,
					}

					switch aV := v.Elt.(type) {
					case *ast.StarExpr:
						fmt.Println("Array pointer value", aV.X.(*ast.Ident).Name)
						fieldData.FieldType = "*" + aV.X.(*ast.Ident).Name
					case *ast.Ident:
						fmt.Println("Array ident", aV.Name)
						fieldData.FieldType = aV.Name
					default:
						panic(fmt.Sprintf("Unsupported array element type: %s", reflect.TypeOf(aV)))
					}

					d.Fields = append(d.Fields, fieldData)
				default:
					panic(fmt.Sprintf("Unsupported type: %s", reflect.TypeOf(v)))
				}
			}
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		panic(err)
	}

	fmt.Println("RENDERED TEMPLATE")
	fmt.Println(buf.String())
}

// x = Parent Object we are unmarshalling
// v = fastJson "Value" that was parsed

func Render(tmpStr string, data interface{}) string {
	tmpl := template.Must(
		template.New("x").
			Funcs(template.FuncMap{
				"lowerCaseFirst": lowerCaseFirst,
				"indentTabs":     indentTabs,
			}).
			Parse(tmpStr))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}

type UnmarshalStringTemplateData struct {
	FieldName string
}

func (s *UnmarshalStringTemplateData) Render() string {
	tmplStr := `x.{{.FieldName}} = string(v.GetStringBytes("{{lowerCaseFirst .FieldName}}"))`
	return Render(tmplStr, s)
}

type UnmarshalPointerObjTemplateData struct {
	FieldName string
	FieldType string
}

func (s *UnmarshalPointerObjTemplateData) Render() string {
	tmplStr := `{
	xField := &{{.FieldType}}{}
	err := xField.UnmarshalFromObj(v.GetObject("{{lowerCaseFirst .FieldName}}"))
	if err != nil {
		return err
	}
	x.{{.FieldName}} = xField
}
`
	return Render(tmplStr, s)
}

type UnmarshalArrayTemplateData struct {
	FieldName string
	FieldType string
}

func (s *UnmarshalArrayTemplateData) Render() string {
	tmplStr := `{
	for _, jV := range v.GetArray("{{lowerCaseFirst .FieldName}}") {
		{{ .UnmarshalElementTemplate }}
		x.{{.FieldName}} = append(x.{{.FieldName}}, y)
	}
}
`
	return Render(tmplStr, s)
}

type UnmarshalMapTemplateData struct {
	FieldName    string
	MapValueType string
}

func (s *UnmarshalMapTemplateData) Render() string {
	tmplStr := `{
	xField := make(map[string]{{.MapValueType}})
	jO := v.GetObject("{{lowerCaseFirst .FieldName}}")
	if jO != nil {
		jO.Visit(func(key []byte, val *fastjson.Value) {
			xField[string(key)] = val.({{.MapValueType}})
		})
	}
	x.{{.FieldName}} = xField
}
`
	return Render(tmplStr, s)
}

type field interface {
	Render() string
}

type fileTemplateData struct {
	StructTypeName string
	Fields         []field
}

const fileTemplateSrc = `
package v2

import (
	"fmt"
	"github.com/valyala/fastjson"
)

{{- $lowerTypeName := lowerCaseFirst .StructTypeName}}

var {{$lowerTypeName}}JsonParserPool fastjson.ParserPool

// Unmarshal unmarshals JSON string to {{.StructTypeName}}.
func (x *{{.StructTypeName}}) Unmarshal(j string) error {
	p := {{$lowerTypeName}}JsonParserPool.Get()
	defer {{$lowerTypeName}}JsonParserPool.Put(p)
	v, err := p.Parse(j)

	if err != nil {
		return fmt.Errorf("failed to parse {{.StructTypeName}} json: %w", err)
	}
    {{range .Fields}}
	{{.Render}}

	{{- end}}

	return nil
}
`
