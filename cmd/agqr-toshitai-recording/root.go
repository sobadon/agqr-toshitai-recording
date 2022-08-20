package main

import (
	"log"

	recbackup "github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/rec-backup"
	"github.com/sobadon/agqr-toshitai-recording/cmd/agqr-toshitai-recording/update"
	"github.com/spf13/cobra"
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

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("%+v", err)
	}
}
