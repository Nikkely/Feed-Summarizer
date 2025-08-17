package summarize

import (
	_ "embed"
	"text/template"
)

//go:embed system_prompt.txt
var systemPrompt string

//go:embed user_prompt.tmpl
var userPromptTmplStr string

//go:embed output.tmpl
var outputTmplStr string

var outputTemplate *template.Template
var userPromptTemplate *template.Template

func init() {
	var err error
	userPromptTemplate, err = template.New("user").Parse(userPromptTmplStr)
	if err != nil {
		panic(err)
	}
	outputTemplate, err = template.New("output").Parse(outputTmplStr)
	if err != nil {
		panic(err)
	}
}
