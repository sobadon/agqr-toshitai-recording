package sqlite

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/sobadon/agqr-toshitai-recording/domain/model/program"
	"github.com/sobadon/agqr-toshitai-recording/domain/repository"
	"github.com/sobadon/agqr-toshitai-recording/internal/errutil"
)

type programSqlite struct {
	ID     int       `db:"id"`
	Title  string    `db:"title"`
	Start  time.Time `db:"start"`
	End    time.Time `db:"end"`
	Status string    `db:"status"`
}

func programSqliteToModelProgram(pgramSqlite programSqlite) program.Program {
	return program.Program{
		ID:     pgramSqlite.ID,
		Title:  pgramSqlite.Title,
		Start:  pgramSqlite.Start,
		End:    pgramSqlite.End,
		Status: program.Status(pgramSqlite.Status),
	}
}

func modelProgramToProgramSqlite(pgram program.Program) programSqlite {
	return programSqlite{
		ID:     pgram.ID,
		Title:  pgram.Title,
		Start:  pgram.Start,
		End:    pgram.End,
		Status: pgram.Status.String(),
	}
}

func NewDB(dbPath string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(errutil.ErrDatabaseOpen, err.Error())
	}
	return db, nil
}

// テーブル作成
func Setup(db *sqlx.DB) error {
	// state: enum('scheduled', 'recording', 'done', 'failed')
	_, err := db.Exec(`create table if not exists programs (
		id integer primary key,
		title text not null,
		start timestamp not null,
		end timestamp not null,
		status text not null,
		created_at timestamp not null default (datetime('now', 'localtime')),
		updated_at timestamp not null default (datetime('now', 'localtime'))
	);`)
	if err != nil {
		return errors.Wrap(errutil.ErrDatabaseQuery, err.Error())
	}

	_, err = db.Exec(`CREATE TRIGGER trigger_updated_at AFTER UPDATE ON programs
		BEGIN
			UPDATE programs SET updated_at = DATETIME('now', 'localtime') WHERE rowid == NEW.rowid;
		END;
		`)
	if err != nil {
		return errors.Wrap(errutil.ErrDatabaseQuery, err.Error())
	}
	return nil
}

type programDatabase struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) repository.ProgramPersistence {
	return &programDatabase{
		DB: db,
	}
}

func (p *programDatabase) Save(ctx context.Context, pgram program.Program) error {
	rows, err := p.DB.QueryxContext(ctx, "select count(*) from programs where id = ?", pgram.ID)
	if err != nil {
		return errors.Wrap(errutil.ErrDatabaseQuery, err.Error())
	}

	var lineCount int
	for rows.Next() {
		err := rows.Scan(&lineCount)
		if err != nil {
			return errors.Wrap(errutil.ErrDatabaseScan, err.Error())
		}
	}

	// 既に番組情報が登録されていれば追加しない
	// TODO: 番組表の変更に対応できない問題がある
	if lineCount != 0 {
		return nil
	}

	pgramSqlite := modelProgramToProgramSqlite(pgram)
	_, err = p.DB.NamedExecContext(ctx, "insert into programs (id, title, start, end, status) values (:id, :title, :start, :end, :status)", pgramSqlite)
	if err != nil {
		return errors.Wrap(errutil.ErrDatabaseQuery, err.Error())
	}
	return nil
}

func (p *programDatabase) LoadStartIn(ctx context.Context, now time.Time, duration time.Duration) ([]program.Program, error) {
	afterAbsoluteTime := now.Add(duration)

	stmt, err := p.DB.PrepareNamedContext(ctx, `select id, title, start, end, status from programs where status = 'scheduled' and :now < start and start < :after`)
	if err != nil {
		return nil, errors.Wrap(errutil.ErrDatabasePrepare, err.Error())
	}

	var pgramsSqlite []programSqlite
	err = stmt.SelectContext(ctx, &pgramsSqlite, map[string]interface{}{"now": now, "after": afterAbsoluteTime})
	if err != nil {
		return nil, errors.Wrap(errutil.ErrDatabaseQuery, err.Error())
	}

	var pgrams []program.Program
	for _, pgramSqlite := range pgramsSqlite {
		pgram := programSqliteToModelProgram(pgramSqlite)
		pgrams = append(pgrams, pgram)
	}

	return pgrams, nil
}

func (p *programDatabase) ChangeStatus(ctx context.Context, pgram program.Program, newStatus program.Status) error {
	_, err := p.DB.NamedExecContext(ctx, `update programs set status = :status where id = :id`, map[string]interface{}{"status": newStatus, "id": pgram.ID})
	if err != nil {
		return errors.Wrap(errutil.ErrDatabaseQuery, err.Error())
	}

	return nil
}
