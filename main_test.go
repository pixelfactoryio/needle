package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mainmocks "go.pixelfactory.io/needle/mocks/main"
)

func TestMainHelp(t *testing.T) {
	is := require.New(t)

	oldArgs := os.Args
	t.Cleanup(func() {
		os.Args = oldArgs
	})
	os.Args = []string{"needle", "--help"}

	is.NotPanics(func() {
		main()
	})
}

func TestMainExecuteError(t *testing.T) {
	is := require.New(t)

	originalExecute := executeCmd
	originalNewLogger := newLogger
	originalExit := osExit

	t.Cleanup(func() {
		executeCmd = originalExecute
		newLogger = originalNewLogger
		osExit = originalExit
	})

	exitCode := 0
	logger := mainmocks.NewErrorLogger(t)
	logger.On("Sync").Return(errors.New("sync failure"))
	logger.On("Error", mock.Anything, mock.Anything).Return().Twice()

	executeCmd = func() error {
		return errors.New("boom")
	}

	newLogger = func() ErrorLogger {
		return logger
	}

	osExit = func(code int) {
		exitCode = code
	}

	main()

	is.Equal(1, exitCode)
	logger.AssertExpectations(t)
}
