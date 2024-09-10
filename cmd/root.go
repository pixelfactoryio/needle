// Package cmd provides root command and subcommands.
package cmd

import (
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.pixelfactory.io/pkg/version"
)

var envPrefix = "NEEDLE"

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	needleCmd, err := NewNeedleCmd()
	if err != nil {
		return err
	}

	cobra.OnInitialize(initConfig)
	return needleCmd.Execute()
}

func initConfig() {
	viper.Set("revision", version.REVISION)
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func bindFlag(flag string) error {
	err := viper.BindPFlag(flag, needleCmd.PersistentFlags().Lookup(flag))
	if err != nil {
		return err
	}

	return nil
}

func newStormClient(dbFile string) (*storm.DB, error) {
	client, err := storm.Open(dbFile)
	if err != nil {
		return nil, err
	}
	return client, nil
}
