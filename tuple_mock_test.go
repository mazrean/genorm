package genorm_test

import "github.com/mazrean/genorm"

//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock/tuple_mock.go -package=mock

type Tuple interface {
	genorm.Tuple
}
