package db_sqlite

import (
	"context"
	"fmt"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type sqliteData struct {
	db *sqlx.DB
}

func NewSQLiteData(db *sqlx.DB) *sqliteData {
	return &sqliteData{
		db: db,
	}
}

func (rp *sqliteData) AddNote(ctx context.Context, data *model.NoteEnc) (err error) {
	action := "add note"

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO note
			(groups_id, title, desc, note)
		VALUES
			(?, ?, ?, ?);
	
	`, data.GroupsID, data.Title, data.Desc, data.Note)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddNote,
			err,
		)
	}

	return nil
}
