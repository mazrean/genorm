package genorm

import (
	"fmt"
	"strings"
)

type Ref[T any] struct{}

// DBType represents the type of database being used.
type DBType uint8

const (
	// MySQL represents MySQL or MariaDB databases.
	// Uses '?' as placeholder.
	MySQL DBType = iota
	// PostgreSQL represents PostgreSQL databases.
	// Uses $1, $2, $3, ... as placeholders.
	PostgreSQL
	// SQLite represents SQLite databases.
	// Uses '?' as placeholder.
	SQLite
)

// genormConfig holds the configuration for GenORM operations.
type genormConfig struct {
	dbType DBType
}

// defaultConfig returns the default configuration (MySQL).
func defaultConfig() genormConfig {
	return genormConfig{
		dbType: MySQL,
	}
}

// Option is a function that configures GenORM.
type Option func(*genormConfig)

// WithDBType sets the database type for placeholder formatting.
func WithDBType(dbType DBType) Option {
	return func(c *genormConfig) {
		c.dbType = dbType
	}
}

// placeholder returns the placeholder string for the given argument index.
// For MySQL/SQLite: returns "?"
// For PostgreSQL: returns "$n" where n is the argument index (1-based).
func (c *genormConfig) placeholder(index int) string {
	switch c.dbType {
	case PostgreSQL:
		return fmt.Sprintf("$%d", index)
	case MySQL, SQLite:
		return "?"
	default:
		return "?"
	}
}

// replacePlaceholders replaces '?' placeholders in an expression query with
// the appropriate placeholder format based on the database type.
// This is used for expressions (like IN clauses) that are built independently.
// startIndex is the starting argument index (1-based).
// Returns the converted query and the next available index.
func (c *genormConfig) replacePlaceholders(query string, startIndex int) (string, int) {
	if c.dbType == MySQL || c.dbType == SQLite {
		// Count the number of placeholders for returning the next index
		count := 0
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				count++
			}
		}
		return query, startIndex + count
	}

	// For PostgreSQL, replace ? with $1, $2, $3, ...
	var result strings.Builder
	result.Grow(len(query) + 10) // Pre-allocate with some extra space for $N placeholders
	index := startIndex
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result.WriteString(fmt.Sprintf("$%d", index))
			index++
		} else {
			result.WriteByte(query[i])
		}
	}
	return result.String(), index
}

// queryBuilder helps build SQL queries with proper placeholder formatting
type queryBuilder struct {
	sb       strings.Builder
	argIndex int
	config   genormConfig
}

func (c *genormConfig) newQueryBuilder() *queryBuilder {
	return &queryBuilder{
		argIndex: 1,
		config:   *c,
	}
}

func (qb *queryBuilder) WriteString(s string) error {
	_, err := qb.sb.WriteString(s)
	return err
}

func (qb *queryBuilder) WriteExprQuery(query string) error {
	converted, nextIndex := qb.config.replacePlaceholders(query, qb.argIndex)
	qb.argIndex = nextIndex
	_, err := qb.sb.WriteString(converted)
	return err
}

func (qb *queryBuilder) String() string {
	return qb.sb.String()
}
