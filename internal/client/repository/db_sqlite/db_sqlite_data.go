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

func (rp *sqliteData) GetDataTitles(ctx context.Context, groupID int) (dataTitles [][]byte, err error) {
	action := "get data titles"

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			title
		FROM
			data
		WHERE
			group_id = ?;
	
	`, groupID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	defer rows.Close()
	var dataTitle []byte
	for rows.Next() {
		if err = rows.Scan(&dataTitle); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		dataTitles = append(dataTitles, dataTitle)
	}

	return dataTitles, nil
}

func (rp *sqliteData) GetDataType(ctx context.Context)

func (rp *sqliteData) AddNote(ctx context.Context, data *model.NoteEnc) (err error) {
	action := "add note"

	row := rp.db.QueryRowContext(ctx, `
		
		INSERT INTO data
			(group_id, data_type, title)
		VALUES
			(?, ?, ?)
		RETURNING id;
	
	`, data.GroupID, data.Type, data.Title)

	var dataID int
	err = row.Scan(&dataID)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddNote,
			err,
		)
	}

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO data_note
			(data_id, desc, note)
		VALUES
			(?, ?, ?);
	
	`, dataID, data.Desc, data.Note)
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

func (rp *sqliteData) GetGroupID(ctx context.Context, accountID int, titleGroup string) (groupID int, err error) {
	action := "get group ID"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			id
		FROM
			groups
		WHERE
			account_id = ? AND title = ?;
	
	`, accountID, titleGroup)

	if err = row.Scan(&groupID); err != nil {
		return -1, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddEmailInGroup,
			err,
		)
	}

	return groupID, nil
}
