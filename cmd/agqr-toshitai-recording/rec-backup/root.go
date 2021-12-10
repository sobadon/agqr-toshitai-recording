package recbackup

import (
	"log"
	"os"
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

	workingDirAbs := outBaseAbsPath + "/working"
	util.Mkdir(workingDirAbs)

	archiveDirAbs := outBaseAbsPath + "/archive"
	util.Mkdir(archiveDirAbs)

	// TODO: timezone
	now := time.Now()

	workingAbsPath := agqr.BuildOutPath(workingDirAbs, now)
	archiveAbsPath := agqr.BuildOutPath(archiveDirAbs, now)

	// 1 hour
	durationSec := 60 * 50
	err = agqr.Rec(durationSec, workingAbsPath)
	if err != nil {
		// TODO: とりあえずの処置なので後でどうにかする
		// ffmpeg が "Connection to tcp://fms2.uniqueradio.jp:443 failed: Operation timed out" で終了したときでも録画を継続させる
		log.Printf("%+v\n", err)
	}

	err = os.Rename(workingAbsPath, archiveAbsPath)
	if err != nil {
		// TODO: とりあえずの処置なので後でどうにかする
		// ffmpeg がコケてファイルが生成されていないときでも録画を継続させる
		log.Printf("%+v\n", err)
	}

	return nil
}
