package database

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUUID(t *testing.T) {
	// UUIDを生成
	uuid, err := GenerateUUID()
	assert.NoError(t, err)

	// UUIDの形式を確認 (8-4-4-4-12の形式で36文字)
	assert.Len(t, uuid, 36)
	assert.Contains(t, uuid, "-")
	parts := strings.Split(uuid, "-")
	assert.Len(t, parts, 5)
	assert.Len(t, parts[0], 8)
	assert.Len(t, parts[1], 4)
	assert.Len(t, parts[2], 4)
	assert.Len(t, parts[3], 4)
	assert.Len(t, parts[4], 12)
}

func TestGenerateUUIDs(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "generate zero UUIDs",
			length:  0,
			wantErr: false,
		},
		{
			name:    "generate one UUID",
			length:  1,
			wantErr: false,
		},
		{
			name:    "generate multiple UUIDs",
			length:  5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuids, err := GenerateUUIDs(tt.length)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, uuids, tt.length)

			// 重複がないことを確認
			seen := make(map[string]bool)
			for _, uuid := range uuids {
				assert.False(t, seen[uuid], "duplicate UUID found")
				seen[uuid] = true
			}
		})
	}
}

func TestDynamicEntity_SaveLoad(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		wantErr bool
	}{
		{
			name: "basic types",
			input: map[string]any{
				"string": "test",
				"int":    42,
				"bool":   true,
				"float":  3.14,
			},
			wantErr: false,
		},
		{
			name: "nested map",
			input: map[string]any{
				"nested": map[string]any{
					"key": "value",
				},
			},
			wantErr: false,
		},
		{
			name: "slice values",
			input: map[string]any{
				"strings": []string{"a", "b", "c"},
				"numbers": []int{1, 2, 3},
			},
			wantErr: false,
		},
		{
			name:    "empty entity",
			input:   map[string]any{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create dynamicEntity from input
			entity := dynamicEntity(tt.input)

			// Test Save
			props, err := entity.Save()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, props, len(tt.input))

			// Test Load
			loaded := make(dynamicEntity)
			err = loaded.Load(props)
			assert.NoError(t, err)

			// Verify all properties were preserved
			for _, prop := range props {
				assert.Equal(t, entity[prop.Name], loaded[prop.Name])
			}
		})
	}
}

func TestDynamicEntity_LoadSaveRoundTrip(t *testing.T) {
	// 初期データ
	original := dynamicEntity{
		"string":  "test",
		"int":     42,
		"float":   3.14,
		"bool":    true,
		"strings": []string{"a", "b", "c"},
	}

	// Save
	props, err := original.Save()
	assert.NoError(t, err)

	// Load into a new entity
	loaded := make(dynamicEntity)
	err = loaded.Load(props)
	assert.NoError(t, err)

	// Compare original and loaded values
	assert.Equal(t, original["string"], loaded["string"])
	assert.Equal(t, original["int"], loaded["int"])
	assert.Equal(t, original["float"], loaded["float"])
	assert.Equal(t, original["bool"], loaded["bool"])
	assert.Equal(t, original["strings"], loaded["strings"])
}
