package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.pixelfactory.io/needle/internal/pkg/version"
)

var envPrefix = "NEEDLE"
var logLevel string

var rootCmd = &cobra.Command{
	Use:           "needle",
	Short:         "needle",
	SilenceErrors: true,
	SilenceUsage:  true,
}

// NewRootCmd create new rootCmd
func NewRootCmd() (*cobra.Command, error) {
	rootCmd.PersistentFlags().StringVar(
		&logLevel,
		"log-level",
		"info",
		"Log level (debug, info, warn, error, fatal, panic)",
	)
	if err := bindFlag("log-level"); err != nil {
		return nil, err
	}

	return rootCmd, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	rootCmd, err := NewRootCmd()
	if err != nil {
		return err
	}

	startCmd, err := NewStartCmd()
	if err != nil {
		return err
	}

	rootCmd.AddCommand(startCmd)
	cobra.OnInitialize(initConfig)
	return rootCmd.Execute()
}

func initConfig() {
	viper.Set("revision", version.REVISION)
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func bindFlag(flag string) error {
	err := viper.BindPFlag(flag, startCmd.PersistentFlags().Lookup(flag))
	if err != nil {
		return err
	}
	return nil
}
