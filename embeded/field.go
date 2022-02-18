package genorm

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mazrean/genorm"
	"github.com/mazrean/genorm/statement"
)

type ColumnField[Type genorm.ExprType] struct {
	IsNull bool
	Valid  bool
	Val    Type
}

func (c *ColumnField[Type]) Scan(src interface{}) error {
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

func (bc *ColumnField[Type]) Value() (val Type, err error) {
	if bc.IsNull {
		return val, statement.ErrNullValue
	}

	if bc.Valid {
		return val, statement.ErrEmptyColumn
	}

	return bc.Val, nil
}

func (bc *ColumnField[Type]) iValue() (val any, err error) {
	return bc.Value()
}

type RelationField[BaseTable Table, RefTable Table, _ JoinedTable] struct{}

func (r RelationField[BaseTable, RefTable, JoinedTable]) Join(expr genorm.TypedTableExpr[JoinedTable, bool]) JoinedTable {
	var (
		baseTable   BaseTable
		refTable    RefTable
		joinedTable JoinedTable
	)

	relation, err := newRelation(join, baseTable, refTable, expr)
	if err != nil {
		joinedTable.AddError(err)
		return joinedTable
	}

	joinedTable.SetRelation(relation)

	return joinedTable
}

func (r RelationField[BaseTable, RefTable, JoinedTable]) LeftJoin(expr genorm.TypedTableExpr[JoinedTable, bool]) JoinedTable {
	var (
		baseTable   BaseTable
		refTable    RefTable
		joinedTable JoinedTable
	)

	relation, err := newRelation(leftJoin, baseTable, refTable, expr)
	if err != nil {
		joinedTable.AddError(err)
		return joinedTable
	}

	joinedTable.SetRelation(relation)

	return joinedTable
}

func (r RelationField[BaseTable, RefTable, JoinedTable]) RightJoin(expr genorm.TypedTableExpr[JoinedTable, bool]) JoinedTable {
	var (
		baseTable   BaseTable
		refTable    RefTable
		joinedTable JoinedTable
	)

	relation, err := newRelation(rightJoin, baseTable, refTable, expr)
	if err != nil {
		joinedTable.AddError(err)
		return joinedTable
	}

	joinedTable.SetRelation(relation)

	return joinedTable
}
