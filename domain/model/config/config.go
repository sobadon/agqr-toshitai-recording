package config

import "time"

type Rec struct {
	// 保存先ディレクトリ
	BasePath string

	// duration 後までの番組をすべて録画する
	PrepareAfter time.Duration

	// 録画前後のマージン
	// 録画全体時間 = マージン + 番組時間 + マージン
	Margin time.Duration
}
