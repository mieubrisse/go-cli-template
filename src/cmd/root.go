package cmd

import (
	"os"

	"github.com/owner-replaceme/project-replaceme/cmd/version"
	"github.com/spf13/cobra"
)

const cmdStr = "project-replaceme"

var rootCmd = &cobra.Command{
	Use:   cmdStr,
	Short: "project-replaceme CLI",
	Long:  "project-replaceme is a command-line tool.",
}

func init() {
	rootCmd.AddCommand(version.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
