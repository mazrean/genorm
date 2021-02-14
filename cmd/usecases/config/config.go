package config

// ConfigReader 設定ファイルからの読み取りのinterface
type ConfigReader interface {
	ReadYAML(path string) error
}
