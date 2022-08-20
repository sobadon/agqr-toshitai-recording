package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

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
