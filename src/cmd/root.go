package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   rootCmdStr,
	Short: "project-replaceme CLI",
	Long:  "project-replaceme is a command-line tool.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
