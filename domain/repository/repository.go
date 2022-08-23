//go:generate mockgen -source=$GOFILE -destination ../../testdata/mock/domain/$GOPACKAGE/$GOFILE
package repository

import (
	"context"
	"time"

	"github.com/sobadon/agqr-toshitai-recording/domain/model/date"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
)

type Station interface {
	GetPrograms(ctx context.Context, date date.Date) ([]program.Program, error)
	Rec(ctx context.Context, basePath string, targetPgram program.Program) error
}

type ProgramPersistence interface {
	Save(ctx context.Context, pgram program.Program) error

	// duration 後までに始まる番組を取得
	LoadStartIn(ctx context.Context, now time.Time, duration time.Duration) ([]program.Program, error)

	// pgram の status を newStatus に変更
	ChangeStatus(ctx context.Context, pgram program.Program, newStatus program.Status) error
}
