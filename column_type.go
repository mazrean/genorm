package genorm

import (
	"database/sql"
	"errors"
	"time"
)

type ColumnField interface {
	sql.Scanner
	iValue() (val any, err error)
}

type ColumnType interface {
	bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string | time.Time | []byte
}

type BasicColumn[Type ColumnType] struct {
	IsNull bool
	Valid  bool
	Val    Type
}

func (c *BasicColumn[Type]) Scan(src interface{}) error {
	if src == nil {
		c.Valid = false
		return nil
	}
	c.Valid = true

	var dest any = c.Val
	switch dest.(type) {
	case bool:
		nb := sql.NullBool{}

		err := nb.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = nb.Valid
		dest = nb.Bool
	case int8:
		ni := sql.NullInt16{}

		err := ni.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ni.Valid
		dest = int8(ni.Int16)
	case int16:
		ns := sql.NullInt16{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ns.Valid
		dest = ns.Int16
	case int32:
		ns := sql.NullInt32{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ns.Valid
		dest = ns.Int32
	case int, int64:
		ni := sql.NullInt64{}

		err := ni.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ni.Valid
		dest = ni.Int64
	case byte: // uint8
		nb := sql.NullByte{}

		err := nb.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = nb.Valid
		dest = nb.Byte
	case uint16:
		ns := sql.NullInt32{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ns.Valid
		dest = uint16(ns.Int32)
	case uint32:
		ns := sql.NullInt64{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ns.Valid
		dest = uint32(ns.Int64)
	case uint64:
		ni := sql.NullInt64{}

		err := ni.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ni.Valid
		dest = uint64(ni.Int64)
	case float32:
		nf := sql.NullFloat64{}

		err := nf.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = nf.Valid
		dest = float32(nf.Float64)
	case float64:
		nf := sql.NullFloat64{}

		err := nf.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = nf.Valid
		dest = nf.Float64
	case string:
		ns := sql.NullString{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = ns.Valid
		dest = ns.String
	case time.Time:
		nt := sql.NullTime{}

		err := nt.Scan(src)
		if err != nil {
			return err
		}

		c.Valid = nt.Valid
		dest = nt.Time
	default:
		return errors.New("unsupported type")
	}

	var ok bool
	c.Val, ok = dest.(Type)
	if !ok {
		return errors.New("failed to convert")
	}

	return nil
}

var (
	ErrNullValue   = errors.New("null value")
	ErrEmptyColumn = errors.New("empty column")
)

func (bc *BasicColumn[Type]) Value() (val Type, err error) {
	if bc.IsNull {
		return val, ErrNullValue
	}

	if bc.Valid {
		return val, ErrEmptyColumn
	}

	return bc.Val, nil
}

func (bc *BasicColumn[Type]) iValue() (val any, err error) {
	return bc.Value()
}

type RelationalColumn[Type ColumnType, BaseColumn Column, RefColumn Column, _ JoinedTable] struct {
	BasicColumn[Type]
}

func (r RelationalColumn[_, BaseColumn, RefColumn, JoinedTable]) Join() JoinedTable {
	var (
		baseColumn  BaseColumn
		refColumn   RefColumn
		joinedTable JoinedTable
	)

	joinedTable.SetRelation(baseColumn, refColumn, Join{})

	return joinedTable
}

func (r RelationalColumn[_, BaseColumn, RefColumn, JoinedTable]) LeftJoin() JoinedTable {
	var (
		baseColumn  BaseColumn
		refColumn   RefColumn
		joinedTable JoinedTable
	)

	joinedTable.SetRelation(baseColumn, refColumn, LeftJoin{})

	return joinedTable
}

func (r RelationalColumn[_, BaseColumn, RefColumn, JoinedTable]) RightJoin() JoinedTable {
	var (
		baseColumn  BaseColumn
		refColumn   RefColumn
		joinedTable JoinedTable
	)

	joinedTable.SetRelation(baseColumn, refColumn, RightJoin{})

	return joinedTable
}
