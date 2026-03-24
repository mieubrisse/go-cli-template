package cmd

import (
	"fmt"

	"github.com/owner-replaceme/project-replaceme/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   versionCmdStr,
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
