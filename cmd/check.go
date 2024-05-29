/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/dguihal/nino/pkg/nino"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run a one time check",
	Long: `
`,
	PreRun: setLogFmt,
	Run:    checkEntrypoint,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func checkEntrypoint(cmd *cobra.Command, args []string) {
	nino.Check(cmd.Context())
}
