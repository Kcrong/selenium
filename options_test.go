package selenigo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Kcrong/selenigo"
)

type MockConvertible struct {
	FieldA string
	FieldB int
}

func (m MockConvertible) ToCapabilities() map[string]interface{} {
	return map[string]interface{}{
		"fieldA": m.FieldA,
		"fieldB": m.FieldB,
	}
}

//nolint:funlen // This is a test file.
func TestMarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    interface{}
		expected map[string]interface{}
		name     string
	}{
		{
			name: "Convertible struct",
			input: MockConvertible{
				FieldA: "test",
				FieldB: 42,
			},
			expected: map[string]interface{}{
				"fieldA": "test",
				"fieldB": 42,
			},
		},
		{
			name: "Regular struct",
			input: struct {
				FirstName string
				LastName  string
				Age       int
			}{
				FirstName: "Alice",
				LastName:  "Kim",
				Age:       30,
			},
			expected: map[string]interface{}{
				"firstName": "Alice",
				"lastName":  "Kim",
				"age":       30,
			},
		},
		{
			name: "Pointer to struct",
			input: &struct {
				FirstName string
				Age       int
			}{
				FirstName: "Bob",
				Age:       25,
			},
			expected: map[string]interface{}{
				"firstName": "Bob",
				"age":       25,
			},
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Non-struct input",
			input:    123,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := selenigo.Marshal(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
