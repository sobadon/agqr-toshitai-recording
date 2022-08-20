package usecase

import (
	"context"

	"github.com/sobadon/agqr-toshitai-recording/domain/model/date"
	"github.com/sobadon/agqr-toshitai-recording/domain/repository"
)

type program struct {
	InfraPersistence repository.ProgramPersistence
	Station          repository.Station
}

func NewProgram(
	infraPersistence repository.ProgramPersistence,
	station repository.Station,
) *program {
	return &program{
		InfraPersistence: infraPersistence,
		Station:          station,
	}
}

func (u *program) Update(ctx context.Context) error {
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
