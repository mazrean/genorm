package genorm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string
		dbType      DBType
		input       string
		expected    string
	}{
		{
			description: "PostgreSQL single placeholder",
			dbType:      PostgreSQL,
			input:       "INSERT INTO users (id) VALUES (?)",
			expected:    "INSERT INTO users (id) VALUES ($1)",
		},
		{
			description: "PostgreSQL multiple placeholders",
			dbType:      PostgreSQL,
			input:       "INSERT INTO users (id, name) VALUES (?, ?)",
			expected:    "INSERT INTO users (id, name) VALUES ($1, $2)",
		},
		{
			description: "PostgreSQL IN clause",
			dbType:      PostgreSQL,
			input:       "SELECT * FROM users WHERE id IN (?, ?, ?)",
			expected:    "SELECT * FROM users WHERE id IN ($1, $2, $3)",
		},
		{
			description: "PostgreSQL multiple rows",
			dbType:      PostgreSQL,
			input:       "INSERT INTO users (id) VALUES (?), (?), (?)",
			expected:    "INSERT INTO users (id) VALUES ($1), ($2), ($3)",
		},
		{
			description: "MySQL unchanged",
			dbType:      MySQL,
			input:       "INSERT INTO users (id) VALUES (?)",
			expected:    "INSERT INTO users (id) VALUES (?)",
		},
		{
			description: "SQLite unchanged",
			dbType:      SQLite,
			input:       "INSERT INTO users (id) VALUES (?)",
			expected:    "INSERT INTO users (id) VALUES (?)",
		},
		{
			description: "PostgreSQL no placeholders",
			dbType:      PostgreSQL,
			input:       "SELECT * FROM users",
			expected:    "SELECT * FROM users",
		},
		{
			description: "PostgreSQL complex query",
			dbType:      PostgreSQL,
			input:       "UPDATE users SET name = ?, email = ? WHERE id = ?",
			expected:    "UPDATE users SET name = $1, email = $2 WHERE id = $3",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			config := genormConfig{dbType: test.dbType}
			result := config.formatQuery(test.input)

			assert.Equal(t, test.expected, result)
		})
	}
}

func TestInLitPlaceholders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description   string
		numLiterals   int
		expectedQuery string
	}{
		{
			description:   "Single literal",
			numLiterals:   1,
			expectedQuery: "(column_name IN (?))",
		},
		{
			description:   "Three literals",
			numLiterals:   3,
			expectedQuery: "(column_name IN (?, ?, ?))",
		},
		{
			description:   "Five literals",
			numLiterals:   5,
			expectedQuery: "(column_name IN (?, ?, ?, ?, ?))",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()

			// Create a simple mock expression
			expr := &ExprStruct[Table, WrappedPrimitive[int]]{
				query: "column_name",
				args:  []ExprType{},
			}

			// Create literals
			literals := make([]WrappedPrimitive[int], test.numLiterals)
			for i := 0; i < test.numLiterals; i++ {
				literals[i] = Wrap(i + 1)
			}

			// Call InLit
			inExpr := InLit(expr, literals...)
			query, args, errs := inExpr.Expr()

			assert.Empty(t, errs)
			assert.Equal(t, test.expectedQuery, query)
			assert.Len(t, args, test.numLiterals)

			// Verify it can be formatted for PostgreSQL
			config := genormConfig{dbType: PostgreSQL}
			pgQuery := config.formatQuery(query)

			// Count placeholders in result
			placeholderCount := strings.Count(pgQuery, "$")
			assert.Equal(t, test.numLiterals, placeholderCount)
		})
	}
}
