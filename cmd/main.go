package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mazrean/gopendb-generator/cmd/infrastructure"
)

func main() {
	cmd, err := infrastructure.InjectRootCmd()
	if err != nil {
		panic(fmt.Errorf("failed to inject: %w", err))
	}

	ctx := context.Background()

	err = cmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
