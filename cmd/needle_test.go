package cmd

import (
	"crypto/tls"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/asdine/storm/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.pixelfactory.io/pkg/server"

	"go.pixelfactory.io/needle/internal/infra/coredns"
	mockscmd "go.pixelfactory.io/needle/mocks/cmd"
)

func resetNeedleCmd() {
	needleCmd = &cobra.Command{
		Use:   "needle",
		Short: "needle",
		RunE:  start,
	}
}

func Test_NewNeedleCmd(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	cmd, err := NewNeedleCmd()
	is.NoError(err)
	is.NotNil(cmd)

	is.NotNil(cmd.PersistentFlags().Lookup("log-level"))
	is.NotNil(cmd.PersistentFlags().Lookup("ca"))
	is.NotNil(cmd.PersistentFlags().Lookup("http-port"))
	is.NotNil(cmd.PersistentFlags().Lookup("https-port"))
}

func Test_ExecuteHelp(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	oldArgs := os.Args
	t.Cleanup(func() {
		os.Args = oldArgs
	})
	os.Args = []string{"needle", "--help"}

	err := Execute()
	is.NoError(err)
}

func Test_Start_LoadX509KeyPairError(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	originalLoad := loadX509KeyPair
	t.Cleanup(func() {
		loadX509KeyPair = originalLoad
	})

	loadX509KeyPair = func(_, _ string) (tls.Certificate, error) {
		return tls.Certificate{}, errors.New("load error")
	}

	err := start(nil, nil)
	is.Error(err)
}

func Test_Start_NewStormClientError(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	originalLoad := loadX509KeyPair
	originalClient := newStormClientFunc
	t.Cleanup(func() {
		loadX509KeyPair = originalLoad
		newStormClientFunc = originalClient
	})

	loadX509KeyPair = func(_, _ string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	newStormClientFunc = func(_ string) (*storm.DB, error) {
		return nil, errors.New("db error")
	}

	err := start(nil, nil)
	is.Error(err)
}

func Test_Start_SuccessWithCoreDNS(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	originalLoad := loadX509KeyPair
	originalClient := newStormClientFunc
	originalServer := newServerFunc
	originalCoreDNS := newCoreDNSServerFunc
	originalCoreDNSEnabled := corednsEnabled
	t.Cleanup(func() {
		loadX509KeyPair = originalLoad
		newStormClientFunc = originalClient
		newServerFunc = originalServer
		newCoreDNSServerFunc = originalCoreDNS
		corednsEnabled = originalCoreDNSEnabled
	})

	loadX509KeyPair = func(_, _ string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	newStormClientFunc = func(_ string) (*storm.DB, error) {
		dbPath := filepath.Join(t.TempDir(), "cache.db")
		return storm.Open(dbPath)
	}

	tlsServer := mockscmd.NewServerRunner(t)
	httpServer := mockscmd.NewServerRunner(t)
	coreDNSServer := mockscmd.NewCoreDNSServer(t)

	tlsServer.On("ListenAndServe").Return(nil)
	httpServer.On("ListenAndServe").Return(nil)
	coreDNSServer.On("Run").Return(nil)

	callCount := 0
	newServerFunc = func(_ ...server.Option) (ServerRunner, error) {
		callCount++
		if callCount == 1 {
			return tlsServer, nil
		}
		return httpServer, nil
	}

	newCoreDNSServerFunc = func(_ ...coredns.Option) CoreDNSServer {
		return coreDNSServer
	}

	corednsEnabled = true

	err := start(nil, nil)
	is.NoError(err)

	mock.AssertExpectationsForObjects(t, tlsServer, httpServer, coreDNSServer)
}

func Test_Start_TLSServerCreateError(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	originalLoad := loadX509KeyPair
	originalClient := newStormClientFunc
	originalServer := newServerFunc
	t.Cleanup(func() {
		loadX509KeyPair = originalLoad
		newStormClientFunc = originalClient
		newServerFunc = originalServer
	})

	loadX509KeyPair = func(_, _ string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	newStormClientFunc = func(_ string) (*storm.DB, error) {
		dbPath := filepath.Join(t.TempDir(), "cache.db")
		return storm.Open(dbPath)
	}

	newServerFunc = func(_ ...server.Option) (ServerRunner, error) {
		return nil, errors.New("tls server error")
	}

	err := start(nil, nil)
	is.Error(err)
}

func Test_Start_HTTPServerCreateError(t *testing.T) {
	is := require.New(t)
	resetNeedleCmd()
	viper.Reset()

	originalLoad := loadX509KeyPair
	originalClient := newStormClientFunc
	originalServer := newServerFunc
	t.Cleanup(func() {
		loadX509KeyPair = originalLoad
		newStormClientFunc = originalClient
		newServerFunc = originalServer
	})

	loadX509KeyPair = func(_, _ string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	newStormClientFunc = func(_ string) (*storm.DB, error) {
		dbPath := filepath.Join(t.TempDir(), "cache.db")
		return storm.Open(dbPath)
	}

	firstServer := mockscmd.NewServerRunner(t)
	firstServer.On("ListenAndServe").Return(nil)

	callCount := 0
	newServerFunc = func(_ ...server.Option) (ServerRunner, error) {
		callCount++
		if callCount == 1 {
			return firstServer, nil
		}
		return nil, errors.New("http server error")
	}

	err := start(nil, nil)
	is.Error(err)

	mock.AssertExpectationsForObjects(t, firstServer)
}
