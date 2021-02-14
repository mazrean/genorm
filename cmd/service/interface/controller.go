package interfaces

import (
	"github.com/mazrean/gopendb-generator/service/data"
)

// Controller 設定ファイルからの読み取りのinterface
type Controller interface {
	Scheme() *data.ControllerOutput
}
