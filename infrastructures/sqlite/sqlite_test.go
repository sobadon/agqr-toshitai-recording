package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
	"github.com/sobadon/agqr-toshitai-recording/internal/timeutil"
)

func tempFilename(t testing.TB) string {
	f, err := os.CreateTemp("", "agqr-toshitai-recording-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func Test_programDatabase_Save(t *testing.T) {
	type args struct {
		pgram program.Program
	}
	tests := []struct {
		name    string
		prepare func(db *sqlx.DB) error
		args    args
		wantErr bool
	}{
		{
			name:    "ok",
			prepare: func(db *sqlx.DB) error { return nil },
			args: args{
				pgram: program.Program{
					ID:    514530,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			wantErr: false,
		},
		{
			name: "既に存在している番組を追加しようとしてもエラーにならない",
			prepare: func(db *sqlx.DB) error {
				_, err := db.Exec(`insert into programs (id, title, start, end) values (
					"514530", "鷲崎健のヨルナイト×ヨルナイト", "2022-08-04 00:00:00+09:00", "2022-08-04 00:30:00+09:00"
				)`)
				return err
			},
			args: args{
				pgram: program.Program{
					ID:    514530,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFilename := tempFilename(t)
			defer os.Remove(tempFilename)
			db, err := sqlx.Open("sqlite3", tempFilename)
			if err != nil {
				t.Fatal(err)
			}

			p := &programDatabase{
				DB: db,
			}

			err = Setup(p.DB)
			if err != nil {
				t.Fatal(err)
			}

			err = tt.prepare(p.DB)
			if err != nil {
				t.Fatal(err)
			}

			// とりあえずエラーなければいいや
			if err := p.Save(context.Background(), tt.args.pgram); (err != nil) != tt.wantErr {
				t.Errorf("programDatabase.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_programDatabase_LoadStartIn(t *testing.T) {
	type args struct {
		now      time.Time
		duration time.Duration
	}
	tests := []struct {
		name    string
		prepare func(db *sqlx.DB) error
		args    args
		want    []program.Program
		wantErr bool
	}{
		{
			name: "番組 1 つ取得できる",
			prepare: func(db *sqlx.DB) error {
				_, err := db.Exec(`insert into programs (id, title, start, end) values
					("514529", "テスト番組名", "2022-08-09 23:50:00+09:00", "2022-08-10 00:00:00+09:00"),
					("514530", "鷲崎健のヨルナイト×ヨルナイト", "2022-08-10 00:00:00+09:00", "2022-08-10 00:30:00+09:00")
				`)
				return err
			},
			args: args{
				now:      time.Date(2022, 8, 9, 23, 59, 30, 0, timeutil.LocationJST()),
				duration: 1 * time.Minute,
			},
			want: []program.Program{
				{
					ID:    514530,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 10, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 10, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			wantErr: false,
		},
		{
			// agqr で 2, 3 分間の番組ってないかも
			name: "番組 2 つ取得できる",
			prepare: func(db *sqlx.DB) error {
				_, err := db.Exec(`insert into programs (id, title, start, end) values
				("514528", "テスト番組名1", "2022-08-09 23:50:00+09:00", "2022-08-10 23:59:00+09:00"),
				("514529", "テスト番組名2", "2022-08-09 23:59:00+09:00", "2022-08-10 00:00:00+09:00"),
				("514530", "鷲崎健のヨルナイト×ヨルナイト", "2022-08-10 00:00:00+09:00", "2022-08-10 00:30:00+09:00")
				`)
				return err
			},
			args: args{
				now:      time.Date(2022, 8, 9, 23, 58, 0, 0, timeutil.LocationJST()),
				duration: 5 * time.Minute,
			},
			want: []program.Program{
				{
					ID:    514529,
					Title: "テスト番組名2",
					Start: time.Date(2022, 8, 9, 23, 59, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 10, 0, 0, 0, 0, timeutil.LocationJST()),
				},
				{
					ID:    514530,
					Title: "鷲崎健のヨルナイト×ヨルナイト",
					Start: time.Date(2022, 8, 10, 0, 0, 0, 0, timeutil.LocationJST()),
					End:   time.Date(2022, 8, 10, 0, 30, 0, 0, timeutil.LocationJST()),
				},
			},
			wantErr: false,
		},
		{
			name: "該当番組がなければ nil を返す",
			prepare: func(db *sqlx.DB) error {
				_, err := db.Exec(`insert into programs (id, title, start, end) values
				("514529", "テスト番組名", "2022-08-09 23:59:00+09:00", "2022-08-10 00:00:00+09:00"),
				("514530", "鷲崎健のヨルナイト×ヨルナイト", "2022-08-10 00:00:00+09:00", "2022-08-10 00:30:00+09:00")
				`)
				return err
			},
			args: args{
				now:      time.Date(2022, 8, 10, 0, 10, 0, 0, timeutil.LocationJST()),
				duration: 5 * time.Minute,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFilename := tempFilename(t)
			defer os.Remove(tempFilename)
			db, err := sqlx.Open("sqlite3", tempFilename)
			if err != nil {
				t.Fatal(err)
			}

			p := &programDatabase{
				DB: db,
			}

			err = Setup(p.DB)
			if err != nil {
				t.Fatal(err)
			}

			err = tt.prepare(p.DB)
			if err != nil {
				t.Fatal(err)
			}

			got, err := p.LoadStartIn(context.Background(), tt.args.now, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("programDatabase.LoadStartIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("programDatabase.LoadStartIn() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
