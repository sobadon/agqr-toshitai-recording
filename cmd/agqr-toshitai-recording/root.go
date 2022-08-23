package main

import (
	recbackup "github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/rec-backup"
	"github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/run"
	"github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/update"
	"github.com/sobadon/agqr-toshitai-recording/internal/logutil"
	"github.com/spf13/cobra"
)

var (
	log = logutil.NewLogger()
)

func main() {
	execute()
}

func execute() {
	var rootCmd = &cobra.Command{
		Use:   "agqr-toshitai-recording",
		Short: "Record agqr streaming",
	}

	rootCmd.AddCommand(recbackup.Execute())
	rootCmd.AddCommand(update.Command())
	rootCmd.AddCommand(run.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Error().Msgf("%+v", err)
	}
}
