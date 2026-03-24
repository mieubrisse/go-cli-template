package version

import (
	"fmt"

	"github.com/owner-replaceme/project-replaceme/internal/buildinfo"
	"github.com/spf13/cobra"
)

const CmdStr = "version"

var Cmd = &cobra.Command{
	Use:   CmdStr,
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), buildinfo.Version)
	},
}
