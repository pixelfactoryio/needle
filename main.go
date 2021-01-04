package main

import (
	"os"

	"go.pixelfactory.io/needle/cmd"
	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"

	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/hosts"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/loop"
)

func main() {
	logger := log.New()
	defer logger.Sync()

	if err := cmd.Execute(); err != nil {
		logger.Error("An unexpected error occured", fields.Error(err))
		os.Exit(1)
	}
}
