// Package jsonify provides utilities for extracting JSON structures from text and formatting them using templates.
package jsonify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

// extractJSONArray extracts either a single JSON object or an array of JSON objects from the input text.
// This function can handle both standalone JSON objects and arrays of objects, even when they are
// embedded within other text content.
//
// The function uses a robust regular expression to find valid JSON structures, handling:
// - Single objects: {...}
// - Arrays of objects: [{...}, {...}]
// - Nested objects and arrays
//
// If a single object is found, it is wrapped in an array for consistent handling.
// This enables uniform processing of both single objects and arrays in the calling code.
func extractJSONArray(input string) ([]string, error) {
	// Find all JSON objects in the input
	re := regexp.MustCompile(`(?s)(\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}|\[[^\[\]]*(?:\[[^\[\]]*\][^\[\]]*)*\])`)
	matches := re.FindAllString(input, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no valid JSON object or array found in the input text")
	}

	var result []string
	for _, match := range matches {
		// Try to parse as array first
		var arr []json.RawMessage
		if err := json.Unmarshal([]byte(match), &arr); err == nil {
			// It's an array, convert each element to string
			for _, item := range arr {
				result = append(result, string(item))
			}
		} else if !strings.HasPrefix(match, "[") {
			// If not an array and doesn't start with '[', treat as single object
			result = append(result, match)
		} else {
			// It starts with '[' but isn't a valid array
			return nil, fmt.Errorf("found a JSON array but its content is invalid or malformed")
		}
	}

	return result, nil
}

// ExtractAndFormat extracts JSON structures from the input text and formats them as a single JSON array.
//
// If multiple JSON objects are found in the input, they are combined into a single array.
// For example, if the input contains {"a": 1} and {"b": 2}, the output will be [{"a": 1}, {"b": 2}].
//
// Parameters:
//   - input: The text containing JSON structure(s)
//   - tmpl: The template to apply to the resulting JSON array
//
// Returns:
//   - any: Formatted JSON array
//   - error: An error if extraction, parsing, or template execution fails
func ExtractAndFormat(input string, tmpl *template.Template) ([]any, error) {
	var err error
	// Extract JSON objects
	jsonObjs, err := extractJSONArray(input)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON: %w", err)
	}

	// Parse each JSON object
	var jsonArr []any
	for _, jsonObj := range jsonObjs {
		var parsed any
		if unmarshalErr := json.Unmarshal([]byte(jsonObj), &parsed); unmarshalErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to parse JSON object: %w", unmarshalErr))
			continue
		}

		var buf bytes.Buffer
		if tmplErr := tmpl.Execute(&buf, parsed); tmplErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to execute template: %w", tmplErr))
			continue
		}

		var newJSONObj any
		if unmarshalErr := json.Unmarshal(buf.Bytes(), &newJSONObj); unmarshalErr != nil {
			err = errors.Join(err, fmt.Errorf("failed to parse JSON object from template output: %w", unmarshalErr))
			continue
		}
		jsonArr = append(jsonArr, newJSONObj)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON string: %w", err)
	}

	return jsonArr, nil
}
