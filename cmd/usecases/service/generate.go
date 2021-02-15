package service

import (
	"context"
	"fmt"

	"github.com/mazrean/gopendb-generator/cmd/usecases/config"
)

// Generate コード生成のservice
type Generate struct {
	config.ConfigReader
}

// NewGenerate Generateのコンストラクタ
func NewGenerate(cr config.ConfigReader) *Generate {
	return &Generate{
		ConfigReader: cr,
	}
}

// Service コード生成のservice
func (g *Generate) Service(ctx context.Context, yamlPath string) error {
	err := g.ConfigReader.ReadYAML(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read yaml: %w", err)
	}

	return nil
}
