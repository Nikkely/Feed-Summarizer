package prompt

import (
	"testing"
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
var userPromptFmt = "タイトル:%s, 本文:%s\n"

func TestNewPromptBuilder(t *testing.T) {
	builder := NewPromptBuilder(systemPrompt, userPromptFmt)

	if builder.SystemPrompt != systemPrompt {
		t.Errorf("expected SystemPrompt to be %q, got %q", systemPrompt, builder.SystemPrompt)
	}

	if builder.UserPromptFmt != userPromptFmt {
		t.Errorf("expected UserPromptFmt to be %q, got %q", userPromptFmt, builder.UserPromptFmt)
	}
}

func TestPromptBuilder_Append(t *testing.T) {
	builder := NewPromptBuilder(systemPrompt, userPromptFmt)
	builder.Append("title1", "body1")
	builder.Append("title2", "body2")

	expected := `タイトル:title1, 本文:body1
タイトル:title2, 本文:body2
`
	if builder.userPrompt != expected {
		t.Errorf("expected userPrompt to be %q, got %q", expected, builder.userPrompt)
	}
}

func TestPromptBuilder_Build(t *testing.T) {
	builder := NewPromptBuilder(systemPrompt, userPromptFmt)
	builder.Append("title1", "body1")
	builder.Append("title2", "body2")

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
	if builder.Build() != expected {
		t.Errorf("expected userPrompt to be %q, got %q", expected, builder.Build())
	}
}
