//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/mazrean/gopendb-generator/cmd/interfaces/handler"
	"github.com/spf13/cobra"
)

func InjectRootCmd() (*cobra.Command, error) {
	wire.Build(
		newRootCmd,
		handler.NewGenerate,
	)

	return &cobra.Command{}, nil
}
