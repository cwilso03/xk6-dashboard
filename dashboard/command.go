package dashboard

import (
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.k6.io/k6/lib/consts"
)

func Execute(uiFS fs.FS) {
	opts := new(options)

	if err := buildRootCmd(opts, uiFS).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

const (
	flagHost   = "host"
	flagPort   = "port"
	flagPeriod = "period"
	flagOpen   = "open"
	flagConfig = "config"

	typeMetric = "Metric"
	typePoint  = "Point"
)

func buildRootCmd(opts *options, uiFS fs.FS) *cobra.Command {
	rootCmd := &cobra.Command{ // nolint:exhaustruct
		Use:   "k6",
		Short: "a next-generation load generator",
		Long:  "\n" + consts.Banner(),
	}

	dashboardCmd := &cobra.Command{ // nolint:exhaustruct
		Use:   "dashboard",
		Short: "xk6-dashboard commands",
	}

	rootCmd.AddCommand(dashboardCmd)

	replayCmd := &cobra.Command{ // nolint:exhaustruct
		Use:   "replay file",
		Short: "load the saved JSON results and replay it for the dashboard UI",
		Long:  "The replay command load the saved JSON results and replay it for the dashboard UI.\nThe compressed file will be automatically decompressed if the file extension is .gz",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := replay(opts, uiFS, args[0]); err != nil {
				return err
			}

			done := make(chan os.Signal, 1)

			signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

			<-done

			return nil
		},
	}

	replayCmd.Flags().SortFlags = false

	flags := replayCmd.PersistentFlags()

	opts = new(options)

	flags.StringVar(&opts.Host, flagHost, defaultHost, "Hostname or IP address for HTTP endpoint (default: '', empty, listen on all interfaces)")
	flags.IntVar(&opts.Port, flagPort, defaultPort, "TCP port for HTTP endpoint (default: 5665), example: 8080")
	flags.DurationVar(&opts.Period, flagPeriod, defaultPeriod, "Event emitting frequency (default: `10s`), example: `1m`")
	flags.BoolVar(&opts.Open, flagOpen, defaultOpen, "Open browser window automatically")
	flags.StringVar(&opts.Config, flagConfig, defaultHost, "UI configuration file location (default: '.dashboard.js')")

	dashboardCmd.AddCommand(replayCmd)

	return rootCmd
}
