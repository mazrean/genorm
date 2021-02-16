package config

//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

// ConfigReader 設定ファイルからの読み取りのinterface
type ConfigReader interface {
	ReadYAML(path string) error
}
