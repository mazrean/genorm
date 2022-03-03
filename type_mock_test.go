package genorm_test

import (
	"github.com/mazrean/genorm"
)

//go:generate go run -workfile=off github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock/${GOFILE} -package=mock

type ExprType interface {
	genorm.ExprType
}

type ColumnFieldExprType interface {
	genorm.ColumnFieldExprType
}
