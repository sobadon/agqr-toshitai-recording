package program

import (
	"context"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/sobadon/agqr-toshitai-recording/infrastructures/agqr"
	"github.com/sobadon/agqr-toshitai-recording/infrastructures/sqlite"
	"github.com/sobadon/agqr-toshitai-recording/usecase"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "program",
		Short: "update program",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	return rootCmd
}

func run() error {
	var config config
	err := env.Parse(&config, env.Options{Prefix: "ATR_"})
	if err != nil {
		return err
	}

	db, err := sqlite.NewDB(config.SqlitePath)
	if err != nil {
		return err
	}

	// マイグレーションとか気にしない
	sqlite.Setup(db)

	infraSqlite := sqlite.New(db)
	stationAgqr := agqr.New()
	ucProgram := usecase.NewProgram(infraSqlite, stationAgqr)

	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 5*time.Second)

	err = ucProgram.Update(ctx)
	if err != nil {
		return err
	}

	return nil
}
