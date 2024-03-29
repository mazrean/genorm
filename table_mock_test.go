package genorm_test

import "github.com/mazrean/genorm"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock/table_mock.go -package=mock

type Table interface {
	genorm.Table
}

type BasicTable interface {
	genorm.BasicTable
}

type JoinedTable interface {
	genorm.JoinedTable
}
