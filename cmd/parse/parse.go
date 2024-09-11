package parse

import (
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLevel   string
	inputFiles []string
	outputFile string
	publicIPv4 string
	publicIPv6 string
)

var parseHostCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse hosts file",
	RunE:  start,
}

func bindFlag(flag string) error {
	err := viper.BindPFlag(flag, parseHostCmd.PersistentFlags().Lookup(flag))
	if err != nil {
		return err
	}

	return nil
}

func NewCmd() (*cobra.Command, error) {
	parseHostCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error, fatal, panic)")
	if err := bindFlag("log-level"); err != nil {
		return nil, err
	}

	parseHostCmd.PersistentFlags().StringSliceVar(&inputFiles, "file", []string{}, "Hosts file URL")
	if err := bindFlag("file"); err != nil {
		return nil, err
	}

	parseHostCmd.PersistentFlags().StringVar(&outputFile, "output", "data/hosts", "Output file")
	if err := bindFlag("output"); err != nil {
		return nil, err
	}

	parseHostCmd.PersistentFlags().StringVar(&publicIPv4, "public-ip", "0.0.0.0", "Public IP")
	if err := bindFlag("public-ip"); err != nil {
		return nil, err
	}

	parseHostCmd.PersistentFlags().StringVar(&publicIPv6, "public-ip6", "::", "Public IPv6")
	if err := bindFlag("public-ip6"); err != nil {
		return nil, err
	}

	return parseHostCmd, nil
}

func start(_ *cobra.Command, _ []string) error {
	hosts := []string{}

	for _, file := range inputFiles {
		bs, err := getHostFile(file)
		if err != nil {
			return err
		}

		hosts = append(hosts, Hosts(bs, publicIPv4, publicIPv6)...)
	}

	// Remove duplicates
	slices.Sort(hosts)
	hosts = slices.Compact(hosts)

	// Write to file
	if err := WriteHosts(hosts, outputFile); err != nil {
		return err
	}

	return nil
}
