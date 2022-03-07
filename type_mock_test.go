package genorm_test

import (
	"github.com/mazrean/genorm"
)

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock/type_mock.go -package=mock

type ExprType interface {
	genorm.ExprType
}

type ColumnFieldExprType interface {
	genorm.ColumnFieldExprType
}
