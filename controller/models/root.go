package models

// Root yamlのルート
type Root struct {
	Config *Config  `yaml:"config"`
	Tables []*Table `yaml:"tables"`
}
