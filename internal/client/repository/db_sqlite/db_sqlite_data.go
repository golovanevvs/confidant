package db_sqlite

import (
	"context"
	"fmt"
	"time"

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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return dataTitles, nil
}

func (rp *sqliteData) GetDataTypes(ctx context.Context, groupID int) (dataTypes []string, err error) {
	action := "get data types"

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			data_type
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
	var dataType string
	for rows.Next() {
		if err = rows.Scan(&dataType); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		dataTypes = append(dataTypes, dataType)
	}

	return dataTypes, nil
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
			customerrors.ErrAddPass,
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
			customerrors.ErrAddPass,
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

func (rp *sqliteData) AddCard(ctx context.Context, data model.CardEnc) (err error) {
	action := "add card"

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
			customerrors.ErrAddCard,
			err,
		)
	}

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO data_card
			(data_id, desc, number, date, name, cvc2, pin, bank)
		VALUES
			(?, ?, ?, ?, ?, ?, ?, ?);
	
	`, dataID, data.Desc, data.Number, data.Date, data.Name, data.CVC2, data.PIN, data.Bank)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddCard,
			err,
		)
	}

	return nil
}

func (rp *sqliteData) GetCard(ctx context.Context, dataID int) (data model.CardEnc, err error) {
	action := "get card"

	data.Type = "card"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			data_card.id, data_card.desc, data_card.number, data_card.date, data_card.name, data_card.cvc2, data_card.pin, data_card.bank
		FROM
			data_card
		INNER JOIN
			data
		ON
			data_card.data_id = data.id
		WHERE
			data_id = ?

	`, dataID)

	if err = row.Scan(&data.ID, &data.Desc, &data.Number, &data.Date, &data.Name, &data.CVC2, &data.PIN, &data.Bank); err != nil {
		return model.CardEnc{}, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrGetCard,
			err,
		)
	}

	return data, nil
}

func (rp *sqliteData) AddFile(ctx context.Context, data model.FileEnc) (err error) {
	action := "add file"

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
			customerrors.ErrAddFile,
			err,
		)
	}

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO data_file
			(data_id, desc, filename, filesize, filedate, file)
		VALUES
			(?, ?, ?, ?, ?, ?);
	
	`, dataID, data.Desc, data.Filename, data.Filesize, data.Filedate, data.File)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddFile,
			err,
		)
	}

	return nil
}

func (rp *sqliteData) GetFile(ctx context.Context, dataID int) (data model.FileEnc, err error) {
	action := "get file"

	data.Type = "file"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			data_file.id, data_file.desc, data_file.filename, data_file.filesize, data_file.filedate
		FROM
			data_file
		INNER JOIN
			data
		ON
			data_file.data_id = data.id
		WHERE
			data_id = ?

	`, dataID)

	if err = row.Scan(&data.ID, &data.Desc, &data.Filename, &data.Filesize, &data.Filedate); err != nil {
		return model.FileEnc{}, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrGetFile,
			err,
		)
	}

	return data, nil
}

func (rp *sqliteData) GetFileForSave(ctx context.Context, dataID int) (dataEnc []byte, err error) {
	action := "get file for save"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			file
		FROM
			data_file
		WHERE
			data_id = ?
	
	`, dataID)

	if err = row.Scan(&dataEnc); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrGetFileForSave,
			err,
		)
	}

	return dataEnc, nil
}

func (rp *sqliteData) GetDataIDs(ctx context.Context, groupServerIDs []int) (dataServerIDs []int, dataNoServerIDs []int, err error) {
	action := "get data IDs"

	dataServerIDs = make([]int, 0)
	dataNoServerIDs = make([]int, 0)

	for _, groupServerID := range groupServerIDs {
		rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			id, id_on_server
		FROM
			data
		WHERE
			(
				SELECT
					id
				FROM
					groups
				WHERE
					id_on_server = ?
			);
	
	`, groupServerID)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		defer rows.Close()

		for rows.Next() {
			var dataID, dataServerID int
			if err = rows.Scan(&dataID, &dataServerID); err != nil {
				return nil, nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}

			if dataServerID == -1 {
				dataNoServerIDs = append(dataNoServerIDs, dataID)
			} else {
				dataServerIDs = append(dataServerIDs, dataServerID)
			}
		}

		if err = rows.Err(); err != nil {
			return nil, nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
	}

	return dataServerIDs, dataNoServerIDs, nil
}

func (rp *sqliteData) GetDataDates(ctx context.Context, dataIDs []int) (dataDates map[int]time.Time, err error) {
	action := "get dates by data IDs"

	dataDates = make(map[int]time.Time)

	for _, dataID := range dataIDs {
		row := rp.db.QueryRowContext(ctx, `
	
			SELECT
				created_at
			FROM
				data
			WHERE
				id = ?;
	
		`, dataID)

		var date time.Time
		if err = row.Scan(&date); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		dataDates[dataID] = date
	}

	return dataDates, nil
}

func (rp *sqliteData) SaveDatas(ctx context.Context, datas []model.Data) (err error) {
	action := "save datas"

	for _, data := range datas {
		row := rp.db.QueryRowContext(ctx, `
		
			SELECT
				id
			FROM
				groups
			WHERE
				id_on_server = ?;
		
		`, data.GroupID)

		var groupID int
		if err = row.Scan(&groupID); err != nil {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		row2 := rp.db.QueryRowContext(ctx, `
		
			INSERT INTO data
				(id_on_server, group_id, data_type, title, created_at)
			VALUES
				(?, ?, ?, ?, ?)
			RETURNING
				id;
		
		`, data.IDOnServer, groupID, data.DataType, data.Title, data.CreatedAt)

		var dataID int
		if err = row2.Scan(&dataID); err != nil {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		switch data.DataType {
		case "note":
			_, err = rp.db.ExecContext(ctx, `
			
				INSERT INTO data_note
					(data_id, desc, note)
				VALUES
					(?, ?, ?);
			
			`, dataID, data.Note.Desc, data.Note.Note)
			if err != nil {
				return fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		case "pass":
			_, err = rp.db.ExecContext(ctx, `
			
				INSERT INTO data_pass
					(data_id, desc, login, pass)
				VALUES
					(?, ?, ?, ?);
			
			`, dataID, data.Pass.Desc, data.Pass.Login, data.Pass.Pass)
			if err != nil {
				return fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		case "card":
			_, err = rp.db.ExecContext(ctx, `
			
				INSERT INTO data_card
					(data_id, desc, number, date, name, cvc2, pin, bank)
				VALUES
					(?, ?, ?, ?, ?, ?, ?, ?);
			
			`, dataID, data.Card.Desc, data.Card.Number, data.Card.Date, data.Card.Name, data.Card.CVC2, data.Card.PIN, data.Card.Bank)
			if err != nil {
				return fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		case "file":
			_, err = rp.db.ExecContext(ctx, `
			
				INSERT INTO data_file
					(data_id, desc, filename, filesize, filedate)
				VALUES
					(?, ?, ?, ?, ?);
			
			`, dataID, data.File.Desc, data.File.Filename, data.File.Filesize, data.File.Filedate)
			if err != nil {
				return fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		}
	}

	return nil
}

func (rp *sqliteData) SaveDataFile(ctx context.Context, dataID int, file []byte) (err error) {
	action := "save file"

	_, err = rp.db.ExecContext(ctx, `
		
		UPDATE
			data_file
		SET
			file = ?
		WHERE
			data_id = ?;
		
		`, file, dataID)

	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return nil
}

func (rp *sqliteData) GetDatasByIDs(ctx context.Context, dataIDs []int) (datas []model.Data, err error) {
	action := "get datas by data IDs"

	datas = make([]model.Data, 0)

	for _, dataID := range dataIDs {
		row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			groups.id_on_server, data.data_type, data.title, data.created_at
		FROM
			data
		INNER JOIN
			groups
		ON
			data.group_id = groups.id
		WHERE
			data.id = ?;

	`, dataID)

		note := model.NoteEnc{}
		pass := model.PassEnc{}
		card := model.CardEnc{}
		file := model.FileEnc{}

		data := model.Data{
			ID: dataID,
		}

		if err = row.Scan(&data.GroupID, &data.DataType, &data.Title, &data.CreatedAt); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		switch data.DataType {
		case "note":
			row2 := rp.db.QueryRowContext(ctx, `
			
				SELECT
					desc, note
				FROM
					data_note
				WHERE
					data_id = ?;
			
			`, dataID)

			if err = row2.Scan(&note.Desc, &note.Note); err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}

		case "pass":
			row2 := rp.db.QueryRowContext(ctx, `
			
				SELECT
					desc, login, pass
				FROM
					data_pass
				WHERE
					data_id = ?;
			
			`, dataID)

			if err = row2.Scan(&pass.Desc, &pass.Login, &pass.Pass); err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}

		case "card":
			row2 := rp.db.QueryRowContext(ctx, `
			
				SELECT
					desc, number, date, name, cvc2, pin, bank
				FROM
					data_card
				WHERE
					data_id = ?;
			
			`, dataID)

			if err = row2.Scan(&card.Desc, &card.Number, &card.Date, &card.Name, &card.CVC2, &card.PIN, &card.Bank); err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}

		case "file":
			row2 := rp.db.QueryRowContext(ctx, `
			
				SELECT
					desc, filename, filesize, filedate
				FROM
					data_file
				WHERE
					data_id = ?;
			
			`, dataID)

			if err = row2.Scan(&file.Desc, &file.Filename, &file.Filesize, &file.Filedate); err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		}

		data.Note = note
		data.Pass = pass
		data.Card = card
		data.File = file

		datas = append(datas, data)
	}

	return datas, nil
}

func (rp *sqliteData) UpdateDataIDsOnServer(ctx context.Context, newDataIDs map[int]int) (err error) {
	action := "update data IDsOnServer"

	for dataID, dataIDOnServer := range newDataIDs {
		_, err = rp.db.ExecContext(ctx, `

		UPDATE
			data
		SET
			id_on_server = ?
		WHERE
			id = ?;

	`, dataIDOnServer, dataID)
	}

	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w",
			customerrors.DBErr,
			action,
			err,
		)
	}

	return nil
}

func (rp *sqliteData) GetDataFile(ctx context.Context, dataID int) (idOnServer int, file []byte, err error) {
	action := "get file by data ID"

	row := rp.db.QueryRowContext(ctx, `
		
		SELECT
			data.id_on_server, data_file.file
		FROM
			data_file
		INNER JOIN
			data
		ON
			data_file.data_id = data.id
		WHERE
			data_file.data_id = ?;
		
	`, dataID)

	if err = row.Scan(&idOnServer, &file); err != nil {
		return -1, nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return idOnServer, file, nil
}
