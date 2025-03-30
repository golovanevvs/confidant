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

func (rp *sqliteData) GetDataIDAndType(ctx context.Context, groupID int, dataTitle string) (dataID int, dataType string, err error) {
	action := "get data ID and data type"

	row := rp.db.QueryRowContext(ctx, `

		SELECT
			id, data_type
		FROM
			data
		WHERE
			group_id = ? AND title = ?;

	`, groupID, dataTitle)

	if err = row.Scan(&dataID, &dataType); err != nil {
		return -1, "", fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return dataID, dataType, nil
}

func (rp *sqliteData) AddNote(ctx context.Context, data model.NoteEnc) (err error) {
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

func (rp *sqliteData) GetNote(ctx context.Context, dataID int) (data model.NoteEnc, err error) {
	action := "get note"

	data.Type = "note"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			data_note.id, data_note.desc, data_note.note
		FROM
			data_note
		INNER JOIN
			data
		ON
			data_note.data_id = data.id
		WHERE
			data_id = ?

	`, dataID)

	if err = row.Scan(&data.ID, &data.Desc, &data.Note); err != nil {
		return model.NoteEnc{}, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrGetNote,
			err,
		)
	}

	return data, nil
}

func (rp *sqliteData) AddPass(ctx context.Context, data model.PassEnc) (err error) {
	action := "add password"

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
			customerrors.ErrAddPAss,
			err,
		)
	}

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO data_pass
			(data_id, desc, login, pass)
		VALUES
			(?, ?, ?, ?);
	
	`, dataID, data.Desc, data.Login, data.Pass)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddPAss,
			err,
		)
	}

	return nil
}

func (rp *sqliteData) GetPass(ctx context.Context, dataID int) (data model.PassEnc, err error) {
	action := "get password"

	data.Type = "pass"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			data_pass.id, data_pass.desc, data_pass.login, data_pass.pass
		FROM
			data_pass
		INNER JOIN
			data
		ON
			data_pass.data_id = data.id
		WHERE
			data_id = ?

	`, dataID)

	if err = row.Scan(&data.ID, &data.Desc, &data.Login, &data.Pass); err != nil {
		return model.PassEnc{}, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrGetPass,
			err,
		)
	}

	return data, nil
}
