package jsonify

import (
	_ "embed"
	"text/template"
)

//go:embed output.json.tmpl
var outputTmplStr string

var OutputTemplate *template.Template

func init() {
	var err error
	OutputTemplate, err = template.New("output").Parse(outputTmplStr)
	if err != nil {
		panic(err)
	}
}
