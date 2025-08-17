package jsonify

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func Test_extractJSONArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
	}{
		{
			name:     "single object",
			input:    "Some text before {\"heading\": \"Test Title\", \"summary\": \"Test Summary\"} and after",
			expected: []string{"{\"heading\": \"Test Title\", \"summary\": \"Test Summary\"}"},
			wantErr:  false,
		},
		{
			name:     "array of objects",
			input:    "Text [{\"heading\": \"Title 1\", \"summary\": \"Sum 1\"}, {\"heading\": \"Title 2\", \"summary\": \"Sum 2\"}]",
			expected: []string{"{\"heading\": \"Title 1\", \"summary\": \"Sum 1\"}", "{\"heading\": \"Title 2\", \"summary\": \"Sum 2\"}"},
			wantErr:  false,
		},
		{
			name:     "no JSON structure",
			input:    "Text with no JSON",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "invalid JSON array",
			input:    "[invalid json]",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractJSONArray(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractAndFormat(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tmplString string
		expected   string
		wantErr    bool
	}{
		{
			name:       "single object",
			input:      "{\"heading\": \"Test\", \"summary\": \"Summary\"}",
			tmplString: "{{range .items}}[{{.heading}}: {{.summary}}]{{end}}",
			expected:   "[Test: Summary]",
			wantErr:    false,
		},
		{
			name:       "array of objects",
			input:      "[{\"heading\": \"H1\", \"summary\": \"S1\"}, {\"heading\": \"H2\", \"summary\": \"S2\"}]",
			tmplString: "{{range .items}}[{{.heading}}: {{.summary}}]{{end}}",
			expected:   "[H1: S1][H2: S2]",
			wantErr:    false,
		},
		{
			name:       "multiple separate objects",
			input:      "{\"heading\": \"H1\", \"summary\": \"S1\"} {\"heading\": \"H2\", \"summary\": \"S2\"}",
			tmplString: "{{range .items}}[{{.heading}}: {{.summary}}]{{end}}",
			expected:   "[H1: S1][H2: S2]",
			wantErr:    false,
		},
		{
			name:       "no JSON structure",
			input:      "Text with no JSON",
			tmplString: "{{.heading}}",
			wantErr:    true,
		},
		{
			name:       "nested objects",
			input:      "[{\"heading\": \"H1\", \"details\": {\"author\": \"A1\"}}, {\"heading\": \"H2\", \"details\": {\"author\": \"A2\"}}]",
			tmplString: "{{range .items}}[{{.heading}} by {{.details.author}}]{{end}}",
			expected:   "[H1 by A1][H2 by A2]",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := template.New("test").Parse(tt.tmplString)
			assert.NoError(t, err, "Template parsing should not fail")

			result, err := ExtractAndFormat(tt.input, tmpl)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
