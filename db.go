package genorm

import "database/sql"

type DB struct {
	db *sql.DB
}

func New(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

func (db *DB) DB() *sql.DB {
	return db.db
}
