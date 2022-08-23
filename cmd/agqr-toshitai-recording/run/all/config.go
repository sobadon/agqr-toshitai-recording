package all

import "time"

type config struct {
	SqlitePath   string        `env:"SQLITE_PATH" envDefault:"db.sqlite3"`
	PrepareAfter time.Duration `env:"PREPARE_DURATION" envDefault:"1m"`
}
