/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dguihal/nino/pkg/nino"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the check command
var daemonizeCmd = &cobra.Command{
	Use:   "daemonize",
	Short: "run nino as a daemon",
	Long: `
`,
	PreRun: setLogFmt,
	Run:    daemonizeEntrypoint,
}

var checkInterval int32

func init() {
	rootCmd.AddCommand(daemonizeCmd)
	daemonizeCmd.PersistentFlags().Int32P("interval", "i", 3600, "Check interval in seconds")

	viper.BindPFlag("interval", daemonizeCmd.PersistentFlags().Lookup("interval"))
	viper.BindEnv("interval", "interval")
}

func daemonizeEntrypoint(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(cmd.Context())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		for {
			select {
			case s := <-signalChan:
				switch s {
				case syscall.SIGHUP:
					//Handle SIGGUP (reload) here
					return
				case os.Interrupt:
					slog.Info("Interrupt received, terminating")
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				slog.Info("Exiting")
				os.Exit(1)
			}
		}
	}()

	//Initialize from vars
	checkInterval = viper.GetInt32("interval")

	if err := daemonRun(ctx); err != nil {
		slog.Error("Runtime failure", slog.Any("error", err))
		os.Exit(1)
	}
}

func daemonRun(ctx context.Context) error {

	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	slog.Info(fmt.Sprintf("Main loop started, check interval defined to %d seconds", checkInterval))

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			slog.Info("Launching a new check")
			nino.Check(ctx)
		}
	}
}
