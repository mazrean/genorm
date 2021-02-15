package service

import "context"

// Generate コード生成のservice
type Generate struct{}

// NewGenerate Generateのコンストラクタ
func NewGenerate() *Generate {
	return &Generate{}
}

// Service コード生成のservice
func (g *Generate) Service(ctx context.Context) error {
	return nil
}
