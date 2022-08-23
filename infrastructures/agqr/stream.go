package agqr

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
	"github.com/sobadon/agqr-toshitai-recording/internal/errutil"
	"github.com/sobadon/agqr-toshitai-recording/internal/fileutil"
)

func (c *client) Rec(ctx context.Context, basePath string, targetPgram program.Program) error {
	file := buildFilepath(basePath, targetPgram)
	err := fileutil.MkdirAllIfNotExist(filepath.Dir(file))
	if err != nil {
		return errors.Wrap(errutil.ErrInternal, err.Error())
	}
	duration := calculateProgramDuration(targetPgram)

	cmd := exec.Command("ffmpeg",
		"-y",
		"-loglevel", "warning", // とりあえず決め打ち
		"-i", c.streamURL.String(),
		"-t", strconv.Itoa(int(duration.Seconds())),
		"-vcodec", "copy",
		"-acodec", "copy",
		file,
	)

	// https://github.com/rs/zerolog/issues/398
	// log.Level(zerolog.InfoLevel).With().Logger() などとしても
	// 出力されるログに loglevel が含まれない
	cmd.Stdout = log.Ctx(ctx).With().Str("level", zerolog.LevelInfoValue).Logger()
	cmd.Stderr = log.Ctx(ctx).With().Str("level", zerolog.LevelWarnValue).Logger()

	log.Ctx(ctx).Debug().Msgf("ffmpeg start ... (program = %+v)", targetPgram)
	log.Ctx(ctx).Debug().Msg(cmd.String())
	err = cmd.Start()
	if err != nil {
		return errors.Wrap(errutil.ErrFfmpeg, err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(errutil.ErrFfmpeg, err.Error())
	}

	return nil
}

func buildFilepath(basePath string, pgram program.Program) string {
	return filepath.Join(basePath, pgram.Start.Format("2006-01-02"), fmt.Sprintf("%s_%s.ts", pgram.Start.Format("2006-01-02_1504"), fileutil.SanitizeReplaceName(pgram.Title)))
}

func calculateProgramDuration(pgram program.Program) time.Duration {
	return pgram.End.Sub(pgram.Start)
}
