package update

import (
	"github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/update/program"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "update",
		Short: "update radio data",
	}

	rootCmd.AddCommand(program.Command())
	return rootCmd
}
