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

// formatQuery replaces all '?' placeholders in the query with the appropriate
// placeholder format based on the database type.
func (c *genormConfig) formatQuery(query string) string {
	if c.dbType == MySQL || c.dbType == SQLite {
		return query
	}

	// For PostgreSQL, replace ? with $1, $2, $3, ...
	var result strings.Builder
	result.Grow(len(query) + 10) // Pre-allocate with some extra space for $N placeholders
	index := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result.WriteString(fmt.Sprintf("$%d", index))
			index++
		} else {
			result.WriteByte(query[i])
		}
	}
	return result.String()
}
