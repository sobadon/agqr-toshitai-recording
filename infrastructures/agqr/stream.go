package agqr

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
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

	// TODO: ffmpeg のログを減らしたり
	cmd := exec.Command("ffmpeg",
		"-y",
		"-allowed_extensions", "ALL",
		"-protocol_whitelist", "file,crypto,http,https,tcp,tls",
		"-i", c.streamURL.String(),
		"-t", strconv.Itoa(int(duration.Seconds())),
		"-vcodec", "copy",
		"-acodec", "copy",
		file,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
