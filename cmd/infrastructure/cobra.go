package infrastructure

import (
	"github.com/mazrean/gopendb-generator/cmd/interfaces/handler"
	"github.com/spf13/cobra"
)

func newRootCmd(generate *handler.Generate) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gopendb",
		Short: "A type safe ORM generator for scheme driven development",
		Long: `Gopendb is a ORM generator for Go.
This application will help you to do type-safe and schema-driven development.`,
	}

	rootCmd.AddCommand(generate.Command)

	return rootCmd
}
