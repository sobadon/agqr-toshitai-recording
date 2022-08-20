package agqr

import (
	"net/http"
	"net/url"
	"time"

	"github.com/sobadon/agqr-toshitai-recording/domain/repository"
)

type client struct {
	httpClient     *http.Client
	programBaseURL *url.URL
	streamURL      *url.URL
}

func New() repository.Station {
	programBaseURL, err := url.Parse("https://www.joqr.co.jp/rss/program/json.php?type=ag")
	if err != nil {
		panic(err)
	}

	// 低画質
	// https://www.uniqueradio.jp/agplayer5/player.php から取得されるもの
	streamURL, err := url.Parse("https://icraft.hs.llnwd.net/agqr10/aandg3.m3u8")
	if err != nil {
		panic(err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = 5 * time.Second

	return &client{
		httpClient:     httpClient,
		programBaseURL: programBaseURL,
		streamURL:      streamURL,
	}
}
