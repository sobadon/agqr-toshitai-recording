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
	// iOS Alexa アプリ + mitmproxxy で勝利
	// （Fire TV Stick + mitmproxy は Amazon.com アカウントでログインしなければならず .com では超 A&G+ のスキルが使えないため敗北）
	// iOS の Alexa アプリではあたかも音声のみをダウンロードしている（`iphone3G` みたいなやつからリンクされている `aac`）のではなく、
	// なんとなんと、普通に動画（`ts`）をダウンロードしてその音声のみを端末のスピーカーから流している
	// 640x360
	// https://www.uniqueradio.jp/agplayerf/hls/amznecho.php
	m3u8URL = "https://icraft.hs.llnwd.net/agqr1/iphone3/HLS_Layer1s.m3u8"
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
