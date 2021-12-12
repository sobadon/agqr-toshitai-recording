package agqr

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	// https://www.uniqueradio.jp/agplayer5/hls/mbr-0-cdn.m3u8
	m3u8URL = "https://icraft.hs.llnwd.net/agqr10/aandg1.m3u8"
)

func Rec(durationSec int, outPath string) error {
	cmd := exec.Command("ffmpeg",
		"-y",
		"-allowed_extensions", "ALL",
		"-protocol_whitelist", "file,crypto,http,https,tcp,tls",
		// reconnect には 1 か true を与える必要がある
		"-reconnect", "true",
		"-i", m3u8URL,
		"-t", strconv.Itoa(durationSec),
		"-vcodec", "copy",
		"-acodec", "copy",
		// "-bsf:a", "aac_adtstoasc",
		"-loglevel", "warning",
		outPath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func BuildOutPath(baseDirAbs string, date time.Time) string {
	// TODO: timezone
	const layout = "2006_01_02_150405"
	d := date.Format(layout)

	return fmt.Sprintf(`%s/ag_%s.ts`, baseDirAbs, d)
}
