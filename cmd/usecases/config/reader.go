package config

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

// Reader 設定ファイルからの読み取りのinterface
type Reader interface {
	ReadYAML(path string) error
}
