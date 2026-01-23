// Package main provides entrypoint for needle.
package main

import (
	"os"

	"go.pixelfactory.io/pkg/observability/log"
	"go.pixelfactory.io/pkg/observability/log/fields"
	"go.uber.org/zap/zapcore"

	"go.pixelfactory.io/needle/cmd"

	_ "github.com/coredns/coredns/plugin/cache"
	_ "github.com/coredns/coredns/plugin/forward"
	_ "github.com/coredns/coredns/plugin/hosts"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/loop"
)

type ErrorLogger interface {
	Error(msg string, fields ...zapcore.Field)
	Sync() error
}

var (
	newLogger  = func() ErrorLogger { return log.New() }
	executeCmd = cmd.Execute
	osExit     = os.Exit
)

func main() {
	logger := newLogger()

	if err := executeCmd(); err != nil {
		logger.Error("An unexpected error occurred", fields.Error(err))

		if syncErr := logger.Sync(); syncErr != nil {
			logger.Error("an error occurred while running logger.Sync()", fields.Error(syncErr))
		}

		osExit(1)
	}
}
