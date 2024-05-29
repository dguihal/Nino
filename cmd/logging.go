package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var debug bool
var logFormat string

func setLogFmt(cmd *cobra.Command, args []string) {

	var logger *slog.Logger

	loggerHandlerOptions := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if debug {
		loggerHandlerOptions.Level = slog.LevelDebug
	}

	switch logFormat {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &loggerHandlerOptions))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerHandlerOptions))
	default:
		panic("Error configuring logger, output must be either 'text' or 'json', '" + logFormat + "' provided")
	}

	slog.SetDefault(logger)
}
