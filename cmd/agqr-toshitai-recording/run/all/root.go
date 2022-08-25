package all

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/recorder"
	"github.com/sobadon/agqr-toshitai-recording/infrastructures/agqr"
	"github.com/sobadon/agqr-toshitai-recording/infrastructures/sqlite"
	"github.com/sobadon/agqr-toshitai-recording/internal/errutil"
	"github.com/sobadon/agqr-toshitai-recording/internal/logutil"
	"github.com/sobadon/agqr-toshitai-recording/internal/timeutil"
	"github.com/sobadon/agqr-toshitai-recording/usecase"
	"github.com/spf13/cobra"
)

var (
	log = logutil.NewLogger()
)

func Command() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "all",
		Short: "run all components",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	return rootCmd
}

func run() error {
	log.Info().Msg("start")

	var config config
	err := env.Parse(&config, env.Options{Prefix: "ATR_"})
	if err != nil {
		return err
	}

	db, err := sqlite.NewDB(config.SqlitePath)
	if err != nil {
		return err
	}

	err = sqlite.Setup(db)
	if err != nil {
		return err
	}

	infraSqlite := sqlite.New(db)
	stationAgar := agqr.New()
	ucRecorder := usecase.NewRecorder(infraSqlite, stationAgar)

	ctx := context.Background()

	scheduler := gocron.NewScheduler(timeutil.LocationJST())

	// job: update
	jobUpdate := func(ctx context.Context, job gocron.Job) {
		ctx = logutil.NewLogger().With().
			Int("job_count", job.RunCount()).
			Str("job", "update").
			Logger().WithContext(ctx)
		err := ucRecorder.Update(ctx)
		if err != nil {
			log.Error().Msgf("%+v", err)
		}
	}
	_, err = scheduler.Every(29*time.Minute).DoWithJobDetails(jobUpdate, ctx)
	if err != nil {
		return errors.Wrap(errutil.ErrScheduler, err.Error())
	}

	// job: update
	jobRec := func(ctx context.Context, job gocron.Job) {
		ctx = logutil.NewLogger().With().
			Int("job_count", job.RunCount()).
			Str("job", "rec").
			Logger().WithContext(ctx)

		useDummyProgram := true
		err := ucRecorder.RecPrepare(ctx, recorder.Config{
			ArchiveDir:   config.ArchiveDir,
			PrepareAfter: config.PrepareAfter,
			Margin:       config.Margin,
		}, useDummyProgram, time.Now().In(timeutil.LocationJST()))
		if err != nil {
			log.Error().Msgf("%+v", err)
		}
	}
	_, err = scheduler.Every(30*time.Minute).DoWithJobDetails(jobRec, ctx)
	if err != nil {
		return errors.Wrap(errutil.ErrScheduler, err.Error())
	}

	scheduler.StartAsync()
	scheduler.RunAllWithDelay(10 * time.Second)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info().Msg("Interrupt")
	defer db.Close()

	return nil
}
