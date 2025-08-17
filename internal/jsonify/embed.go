package jsonify

import (
	_ "embed"
	"text/template"
)

//go:embed output.json.tmpl
var outputTmplStr string

// OutputTemplate is a pre-defined template for formatting JSON data.
// It is loaded from output.json.tmpl and defines the structure of JSON objects.
var OutputTemplate *template.Template

func init() {
	var err error
	OutputTemplate, err = template.New("output").Parse(outputTmplStr)
	if err != nil {
		panic(err)
	}
}
