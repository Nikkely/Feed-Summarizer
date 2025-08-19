package prompt

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

var systemPrompt = `
あなたは優秀なニュース記者です。以下に示すRSSの更新情報を元に、わかりやすく、かつ簡潔にニュース記事を作成してください。

- 重要なポイントを押さえ、事実に基づいて正確にまとめてください。
- 専門用語や難しい言葉は避け、一般読者が理解しやすい表現を心がけてください。
- 不必要な冗長表現は避け、要点を簡潔に伝えてください。
- 記事のトーンは客観的かつ中立的で、感情的にならないようにしてください。
- 文章の流れが自然になるように構成してください。

では、以下のRSS情報を基にニュース記事を作成してください。
`
var templateStr = "タイトル:{{.title}}, 本文:{{.body}}\n"

func getTemplate(str string) *template.Template {
	temp, err := template.New("info").Parse(str)
	if err != nil {
		panic("failed to parse user prompt template: " + err.Error())
	}
	return temp
}

func TestNewPromptBuilder(t *testing.T) {
	systemPrompt := "System: "
	validTemplateStr := "User: {{.Name}}\n"
	builder := NewPromptBuilder(systemPrompt, getTemplate(validTemplateStr))
	assert.Equal(t, systemPrompt, builder.SystemPrompt)
	assert.NotNil(t, builder.UserPromptTemplate)
}

func TestPromptBuilder_Append(t *testing.T) {
	builder := NewPromptBuilder(systemPrompt, getTemplate(templateStr))

	builder.Append(map[string]string{"title": "title1", "body": "body1"}).Append(map[string]string{"title": "title2", "body": "body2"})

	expected := `タイトル:title1, 本文:body1
タイトル:title2, 本文:body2
`
	assert.Equal(t, expected, builder.userPrompt)
}

func TestPromptBuilder_Append_Error(t *testing.T) {
	systemPrompt := "System: "
	templateStr := "User: {{.Name}}\n"
	builder := NewPromptBuilder(systemPrompt, getTemplate(templateStr))

	// Pass invalid data type to trigger an error
	invalidVars := 123 // Not a map or struct
	result := builder.Append(invalidVars)
	assert.Nil(t, result)
}

func TestPromptBuilder_Build(t *testing.T) {
	builder := NewPromptBuilder(systemPrompt, getTemplate(templateStr))

	builder.Append(map[string]string{"title": "title1", "body": "body1"}).Append(map[string]string{"title": "title2", "body": "body2"})
	expected := `
あなたは優秀なニュース記者です。以下に示すRSSの更新情報を元に、わかりやすく、かつ簡潔にニュース記事を作成してください。

- 重要なポイントを押さえ、事実に基づいて正確にまとめてください。
- 専門用語や難しい言葉は避け、一般読者が理解しやすい表現を心がけてください。
- 不必要な冗長表現は避け、要点を簡潔に伝えてください。
- 記事のトーンは客観的かつ中立的で、感情的にならないようにしてください。
- 文章の流れが自然になるように構成してください。

では、以下のRSS情報を基にニュース記事を作成してください。
タイトル:title1, 本文:body1
タイトル:title2, 本文:body2
`
	result := builder.Build()
	assert.Equal(t, expected, result)
}
