package codegen

import "testing"

func TestEscapeTag(t *testing.T) {
	tests := []struct {
		description string
		tag         string
		expected    string
	}{
		{
			description: "normal",
			tag:         "hoge",
			expected:    "hoge",
		},
		{
			description: "with double quote",
			tag:         `"piyo"`,
			expected:    `\"piyo\"`,
		},
		{
			description: "with single quote",
			tag:         `'piyo'`,
			expected:    `'piyo'`,
		},
		{
			description: "with back quote",
			tag:         "`piyo`",
			expected:    "`piyo`",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := escapeTag(test.tag)
			if actual != test.expected {
				t.Errorf("expected: %s, actual: %s", test.expected, actual)
			}
		})
	}
}
