package config

import "github.com/mazrean/gopendb-generator/cmd/domain"

// Table 設定ファイルのtableのinterface
type Table interface {
	GetAll() ([]*domain.Table, error)
}
