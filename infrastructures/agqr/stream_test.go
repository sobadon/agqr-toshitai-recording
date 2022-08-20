package agqr

import (
	"testing"
	"time"

	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
	"github.com/sobadon/agqr-toshitai-recording/internal/timeutil"
)

func Test_buildFilepath(t *testing.T) {
	type args struct {
		basePath string
		pgram    program.Program
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basePath の末尾に / が存在しなくてもよい",
			args: args{
				basePath: "/archive",
				pgram: program.Program{
					ID:    514569,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			want: "/archive/2022-08-04/2022-08-04_0000_鷲崎健のヨルナイト×ヨルナイト.ts",
		},
		{
			name: "basePath の末尾に / が存在してもよい",
			args: args{
				basePath: "/archive",
				pgram: program.Program{
					ID:    514569,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			want: "/archive/2022-08-04/2022-08-04_0000_鷲崎健のヨルナイト×ヨルナイト.ts",
		},
		{
			name: "Title にファイル名に使えない or 面倒な文字が含まれていたら使えるものに置換される",
			args: args{
				basePath: "/archive",
				pgram: program.Program{
					ID:    514569,
					Title: "超!A&G+スペシャル",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),  // ダミー
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()), // ダミー
				},
			},
			want: "/archive/2022-08-04/2022-08-04_0000_超！A＆G＋スペシャル.ts",
		},
		{
			name: "日付、時刻が 1 ケタ・1 ケタ",
			args: args{
				basePath: "/archive",
				pgram: program.Program{
					ID:    514569,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			want: "/archive/2022-08-04/2022-08-04_0000_鷲崎健のヨルナイト×ヨルナイト.ts",
		},
		{
			name: "日付、時刻が 2 ケタ・2 ケタ",
			args: args{
				basePath: "/archive",
				pgram: program.Program{
					ID:    514569,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 12, 31, 12, 30, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 12, 31, 13, 0, 0, 0, timeutil.LocationJST()),
				},
			},
			want: "/archive/2022-12-31/2022-12-31_1230_鷲崎健のヨルナイト×ヨルナイト.ts",
		},
		{
			name: "絶対パス",
			args: args{
				basePath: "/archive",
				pgram: program.Program{
					ID:    514569,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			want: "/archive/2022-08-04/2022-08-04_0000_鷲崎健のヨルナイト×ヨルナイト.ts",
		},
		{
			name: "相対パス",
			args: args{
				basePath: "./archive",
				pgram: program.Program{
					ID:    514569,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			want: "archive/2022-08-04/2022-08-04_0000_鷲崎健のヨルナイト×ヨルナイト.ts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFilepath(tt.args.basePath, tt.args.pgram); got != tt.want {
				t.Errorf("buildFilepath() = %v, want %v", got, tt.want)
			}
		})
	}
}
