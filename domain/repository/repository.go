package repository

import (
	"context"

	"github.com/sobadon/agqr-toshitai-recording/domain/model/date"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
)

type Station interface {
	GetPrograms(ctx context.Context, date date.Date) ([]program.Program, error)
}
