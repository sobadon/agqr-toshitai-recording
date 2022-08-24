package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/recorder"
	"github.com/sobadon/agqr-toshitai-recording/internal/errutil"
	"github.com/sobadon/agqr-toshitai-recording/internal/timeutil"
	mock_repository "github.com/sobadon/agqr-toshitai-recording/testdata/mock/domain/repository"
)

func Test_recorder_rec(t *testing.T) {
	pgramNormal := program.Program{
		ID:     514569,
		Title:  "鷲崎健のヨルナイト×ヨルナイト",
		Start:  time.Date(2022, 8, 4, 0, 0, 0, 0, timeutil.LocationJST()),
		End:    time.Date(2022, 8, 4, 0, 30, 0, 0, timeutil.LocationJST()),
		Status: program.StatusScheduled,
	}

	type fields struct {
		InfraPersistence *mock_repository.MockProgramPersistence
		Station          *mock_repository.MockStation
	}
	type args struct {
		config      recorder.Config
		targetPgram program.Program
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
	}{
		// mock の呼び出し回数だけで、問題ない　or 問題ありを判断するものとする
		// この rec() が error を返さないので（返させたくなかったので）、こうなってしまった

		{
			name: "何ら問題なく録画に成功",
			prepare: func(f *fields) {
				f.InfraPersistence.EXPECT().
					ChangeStatus(gomock.Any(), pgramNormal, program.StatusRecording).
					Return(nil)
				f.Station.EXPECT().
					Rec(gomock.Any(), "/archive", pgramNormal).
					Return(nil)
				f.InfraPersistence.EXPECT().
					ChangeStatus(gomock.Any(), pgramNormal, program.StatusDone).
					Return(nil)
			},
			args: args{
				config: recorder.Config{
					BasePath:     "/archive",
					PrepareAfter: 1 * time.Minute,
				},
				targetPgram: pgramNormal,
			},
		},
		{
			name: "一度 ffmpeg が異常終了したとしてもリトライによって録画を継続",
			prepare: func(f *fields) {
				f.InfraPersistence.EXPECT().
					ChangeStatus(gomock.Any(), pgramNormal, program.StatusRecording).
					Return(nil)
				f.Station.EXPECT().
					Rec(gomock.Any(), "/archive", pgramNormal).
					Return(errors.Wrap(errutil.ErrFfmpeg, "something error")).
					Times(1)
				f.Station.EXPECT().
					Rec(gomock.Any(), "/archive", pgramNormal).
					Return(nil).
					Times(1)
				f.InfraPersistence.EXPECT().
					ChangeStatus(gomock.Any(), pgramNormal, program.StatusDone).
					Return(nil)
			},
			args: args{
				config: recorder.Config{
					BasePath:     "/archive",
					PrepareAfter: 1 * time.Minute,
				},
				targetPgram: pgramNormal,
			},
		},
		{
			name: "リトライを最大回数実施したが変わらず異常であるので録画を異常終了させる",
			prepare: func(f *fields) {
				f.InfraPersistence.EXPECT().
					ChangeStatus(gomock.Any(), pgramNormal, program.StatusRecording).
					Return(nil)
				f.Station.EXPECT().
					Rec(gomock.Any(), "/archive", pgramNormal).
					Return(errors.Wrap(errutil.ErrFfmpeg, "something error")).
					Times(4) // retryCount=0, 1, 2, 3 の計 4 回トライする
				f.InfraPersistence.EXPECT().
					ChangeStatus(gomock.Any(), pgramNormal, program.StatusFailed).
					Return(nil)
			},
			args: args{
				config: recorder.Config{
					BasePath:     "/archive",
					PrepareAfter: 1 * time.Minute,
				},
				targetPgram: pgramNormal,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockInfraPersistence := mock_repository.NewMockProgramPersistence(ctrl)
			mockStation := mock_repository.NewMockStation(ctrl)
			u := &ucRecorder{
				InfraPersistence: mockInfraPersistence,
				Station:          mockStation,
			}
			f := &fields{
				InfraPersistence: mockInfraPersistence,
				Station:          mockStation,
			}
			tt.prepare(f)

			u.rec(context.Background(), tt.args.config, tt.args.targetPgram)
		})
	}
}
