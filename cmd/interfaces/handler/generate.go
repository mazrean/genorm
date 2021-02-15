package handler

import (
	"github.com/spf13/cobra"
)

// Generate コード生成コマンド
type Generate struct {
	*cobra.Command
}

// NewGenerate Generateのコンストラクタ
func NewGenerate() *Generate {
	generate := &Generate{}
	generate.Command = &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   "Generate ORM",
		Long:    `Generate ORM`,
		RunE:    generate.Handler,
	}

	return generate
}

// Handler generateのハンドラー
func (*Generate) Handler(cmd *cobra.Command, args []string) error {
	return nil
}
