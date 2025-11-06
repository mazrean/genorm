package genorm

import "fmt"

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

// formatPlaceholder formats a placeholder based on the database type.
// For MySQL/SQLite: returns "?"
// For PostgreSQL: returns "$n" where n is the argument index (1-based).
func (c *genormConfig) formatPlaceholder(index int) string {
	switch c.dbType {
	case PostgreSQL:
		return fmt.Sprintf("$%d", index)
	case MySQL, SQLite:
		return "?"
	default:
		return "?"
	}
}

// formatQuery replaces all '?' placeholders in the query with the appropriate
// placeholder format based on the database type.
func (c *genormConfig) formatQuery(query string) string {
	if c.dbType == MySQL || c.dbType == SQLite {
		return query
	}

	// For PostgreSQL, replace ? with $1, $2, $3, ...
	result := ""
	index := 1
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			result += fmt.Sprintf("$%d", index)
			index++
		} else {
			result += string(query[i])
		}
	}
	return result
}
