package genorm

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

type ExprType interface {
	driver.Valuer
}

type ColumnFieldExprType interface {
	sql.Scanner
	driver.Valuer
}

type ColumnFieldExprTypePointer[T ExprType] interface {
	ColumnFieldExprType
	*T
}

type ExprPrimitive interface {
	bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string | time.Time
}

type WrappedPrimitive[T ExprPrimitive] struct {
	valid bool
	val   T
}

func Wrap[T ExprPrimitive](val T) WrappedPrimitive[T] {
	return WrappedPrimitive[T]{
		valid: true,
		val:   val,
	}
}

func (wp *WrappedPrimitive[T]) Scan(src any) error {
	var dest any = wp.val
	switch dest.(type) {
	case bool:
		nb := sql.NullBool{}

		err := nb.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = nb.Valid
		dest = nb.Bool
	case int8:
		ni := sql.NullInt16{}

		err := ni.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ni.Valid
		dest = int8(ni.Int16)
	case int16:
		ns := sql.NullInt16{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ns.Valid
		dest = ns.Int16
	case int32:
		ns := sql.NullInt32{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ns.Valid
		dest = ns.Int32
	case int, int64:
		ni := sql.NullInt64{}

		err := ni.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ni.Valid
		dest = ni.Int64
	case byte: // uint8
		nb := sql.NullByte{}

		err := nb.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = nb.Valid
		dest = nb.Byte
	case uint16:
		ns := sql.NullInt32{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ns.Valid
		dest = uint16(ns.Int32)
	case uint32:
		ns := sql.NullInt64{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ns.Valid
		dest = uint32(ns.Int64)
	case uint64:
		ni := sql.NullInt64{}

		err := ni.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ni.Valid
		dest = uint64(ni.Int64)
	case float32:
		nf := sql.NullFloat64{}

		err := nf.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = nf.Valid
		dest = float32(nf.Float64)
	case float64:
		nf := sql.NullFloat64{}

		err := nf.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = nf.Valid
		dest = nf.Float64
	case string:
		ns := sql.NullString{}

		err := ns.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = ns.Valid
		dest = ns.String
	case time.Time:
		nt := sql.NullTime{}

		err := nt.Scan(src)
		if err != nil {
			return err
		}

		wp.valid = nt.Valid
		dest = nt.Time
	default:
		return errors.New("unsupported type")
	}

	var ok bool
	wp.val, ok = dest.(T)
	if !ok {
		return errors.New("failed to convert")
	}

	return nil
}

func (wp WrappedPrimitive[_]) Value() (driver.Value, error) {
	if !wp.valid {
		return nil, ErrNullValue
	}

	return wp.val, nil
}

func (wp WrappedPrimitive[T]) Val() (T, bool) {
	return wp.val, wp.valid
}
