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
	config.Table
}

// NewGenerate Generateのコンストラクタ
func NewGenerate(cr config.Reader, cf config.Config, ct config.Table) *Generate {
	return &Generate{
		Reader: cr,
		Config: cf,
		Table:  ct,
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

	tables, err := g.Table.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	return nil
}
