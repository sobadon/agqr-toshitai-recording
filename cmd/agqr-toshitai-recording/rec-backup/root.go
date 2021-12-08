package recbackup

import (
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sobadon/agqr-toshitai-recording/agqr"
	"github.com/sobadon/agqr-toshitai-recording/util"
	"github.com/spf13/cobra"
)

func Execute() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "rec-backup",
		Short: "Record agqr streaming (backup)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return recLoop()
		},
	}

	return rootCmd
}

func recLoop() error {
	for {
		err := rec()
		if err != nil {
			return err
		}
	}
}

func rec() error {
	// temp
	const baseDir = "./rec"

	outBaseAbsPath, err := filepath.Abs(baseDir)
	if err != nil {
		return errors.WithStack(err)
	}

	util.Mkdir(outBaseAbsPath)

	// TODO: timezone
	now := time.Now()

	outAbsPath := agqr.BuildOutPath(outBaseAbsPath, now)

	// 1 hour
	durationSec := 60 * 60
	err = agqr.Rec(durationSec, outAbsPath)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
