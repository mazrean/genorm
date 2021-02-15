package service

import (
	"context"
	"fmt"

	"github.com/mazrean/gopendb-generator/cmd/usecases/config"
)

// Generate コード生成のservice
type Generate struct {
	config.Reader
	config.Config
}

// NewGenerate Generateのコンストラクタ
func NewGenerate(cr config.Reader, cf config.Config) *Generate {
	return &Generate{
		Reader: cr,
		Config: cf,
	}
}

// Service コード生成のservice
func (g *Generate) Service(ctx context.Context, yamlPath string) error {
	err := g.Reader.ReadYAML(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read yaml: %w", err)
	}

	config, err := g.Config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	return nil
}
