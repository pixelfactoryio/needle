// Package main provides entrypoint for needle.
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

	if err := cmd.Execute(); err != nil {
		logger.Error("An unexpected error occurred", fields.Error(err))

		err := logger.Sync()
		if err != nil {
			logger.Error("an error occurred while running logger.Sync()", fields.Error(err))
		}

		os.Exit(1)
	}
}
