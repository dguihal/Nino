package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dguihal/nino/pkg/nino"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// https://ieftimov.com/posts/four-steps-daemonize-your-golang-programs/

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

	if err := viper.BindPFlag("interval", daemonizeCmd.PersistentFlags().Lookup("interval")); err != nil {
		slog.Error("Runtime failure", slog.Any("error", err))
		os.Exit(1)
	}

	if err := viper.BindEnv("interval", "interval"); err != nil {
		slog.Error("Runtime failure", slog.Any("error", err))
		os.Exit(1)
	}
}

func daemonizeEntrypoint(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithCancel(cmd.Context())

	go func() {
		server := &http.Server{
			Addr:         ":2112",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		http.Handle("/metrics", promhttp.Handler())
		slog.Info("Prometeus compatible metrics available on 'http://:2112/metrics'")

		if err := server.ListenAndServe(); err != nil {
			slog.Error("Runtime failure", slog.Any("error", err))
			os.Exit(1)
		}
	}()

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
					// Handle SIGGUP (reload) here
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

	// Initialize from vars
	checkInterval = viper.GetInt32("interval")

	if err := daemonRun(ctx); err != nil {
		slog.Error("Runtime failure", slog.Any("error", err))
		os.Exit(1)
	}
}

func daemonRun(ctx context.Context) error {

	//	slog.Info("Launching an initial check")
	//	nino.Check(ctx)

	ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)
	defer ticker.Stop()
	slog.Info(fmt.Sprintf("Main loop started, check interval defined to %d seconds", checkInterval))

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			slog.Info("###############################################################")
			slog.Info(fmt.Sprintf("Current date and time is: %s", time.Now().String()))
			slog.Info("####################################################################")
			ticker.Stop()
			slog.Info("Launching a new check")
			nino.Check(ctx)
			ticker.Reset(time.Duration(checkInterval) * time.Second)
		}
	}
}
