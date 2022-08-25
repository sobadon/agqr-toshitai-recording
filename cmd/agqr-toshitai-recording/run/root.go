package run

import (
	"github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/run/all"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "run",
		Short: "run component",
	}

	rootCmd.AddCommand(all.Command())
	return rootCmd
}
