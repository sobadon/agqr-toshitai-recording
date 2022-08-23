package usecase

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/date"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
	"github.com/sobadon/agqr-toshitai-recording/domain/repository"
)

type recorder struct {
	InfraPersistence repository.ProgramPersistence
	Station          repository.Station
}

func NewProgram(
	infraPersistence repository.ProgramPersistence,
	station repository.Station,
) *recorder {
	return &recorder{
		InfraPersistence: infraPersistence,
		Station:          station,
	}
}

func (u *recorder) Update(ctx context.Context) error {
	programs, err := u.Station.GetPrograms(ctx, date.NewFromToday())
	if err != nil {
		return err
	}

	for _, program := range programs {
		err := u.InfraPersistence.Save(ctx, program)
		if err != nil {
			return err
		}
	}

	return nil
}

type RecConfig struct {
	// 保存先ディレクトリ
	BasePath string

	// duration 後までの番組をすべて録画する
	PrepareAfter time.Duration
}

// 録画する
// これは一定時間毎に呼び出されなければならない
// Rec という名前、ややこしい
func (u *recorder) RecPrepare(ctx context.Context, config RecConfig, isDebug bool, now time.Time) error {
	var targetPgrams []program.Program
	if isDebug {
		log.Ctx(ctx).Debug().Msg("use dummy programs")
		// ダミー番組の status を変更できていない！
		targetPgrams = program.Dummies(now)
	} else {
		var err error
		targetPgrams, err = u.InfraPersistence.LoadStartIn(ctx, now, config.PrepareAfter)
		if err != nil {
			return err
		}
	}

	for _, targetPgram := range targetPgrams {
		// rec 内部で雑にエラーハンドリングしちゃう
		go u.rec(ctx, config, targetPgram)
	}

	return nil
}

// 録画処理を呼び出す
// 内部でリトライあり
// これは goroutine として呼び出されることを想定
// エラーが発生すれば Error レベルでログを出力してしまう
// Fatal は exit してしまうので使わない
func (u *recorder) rec(ctx context.Context, config RecConfig, targetPgram program.Program) {
	// retryCount=0, 1, 2, 3 の計 4 回トライする
	const retryMaxCount = 3
	retryCount := 0

	err := u.InfraPersistence.ChangeStatus(ctx, targetPgram, program.StatusRecording)
	if err != nil {
		log.Ctx(ctx).Warn().Msgf("fail to change status (scheduled -> recording): %+v", err)
		return
	}

	for retryCount <= retryMaxCount {
		log.Ctx(ctx).Debug().Msgf("rec ... (retryCount = %d)", retryCount)
		err = u.Station.Rec(ctx, config.BasePath, targetPgram)
		if err == nil {
			log.Ctx(ctx).Debug().Msgf("successfully recorded (program.ID = %d)", targetPgram.ID)
			err := u.InfraPersistence.ChangeStatus(ctx, targetPgram, program.StatusDone)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to change status (recording -> done): %+v", err)
				return
			}
			return
		}
		log.Ctx(ctx).Warn().Msgf("fail to rec: %+v", err)
		retryCount++
	}

	// goroutine で呼び出されたとき他の進行中の録画 (goroutine) を犠牲にせずに error を戻し返すのが面倒？
	log.Ctx(ctx).Error().Msgf("retry count exceeded (retryMaxCount = %d)", retryMaxCount)
	err = u.InfraPersistence.ChangeStatus(ctx, targetPgram, program.StatusFailed)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to change status (recording -> failed): %+v", err)
	}
}
