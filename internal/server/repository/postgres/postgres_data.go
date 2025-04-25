package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/jmoiron/sqlx"
)

type postgresData struct {
	db   *sqlx.DB
	rpgp *postgresGroups
}

func NewPostgresData(db *sqlx.DB, rpgp *postgresGroups) *postgresData {
	return &postgresData{
		db:   db,
		rpgp: rpgp,
	}
}

func (rp *postgresData) GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error) {
	action := "get data IDs"

	// getting group IDs
	var groupIDs []int
	groupIDs, err = rp.rpgp.GetGroupIDs(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	if len(groupIDs) > 0 {
		dataIDs = make([]int, 0)
		for _, groupID := range groupIDs {
			rows, err := rp.db.QueryContext(ctx, `
			
				SELECT
					id
				FROM
					data
				WHERE
					group_id = $1;
	
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

			for rows.Next() {
				var dataID int
				if err = rows.Scan(&dataID); err != nil {
					return nil, fmt.Errorf(
						"%s: %s: %w: %w",
						customerrors.DBErr,
						action,
						customerrors.ErrDBInternalError500,
						err,
					)
				}
				dataIDs = append(dataIDs, dataID)
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
		}
		return dataIDs, nil
	} else {
		return nil, nil
	}
}

func (rp *postgresData) GetDataDates(ctx context.Context, dataIDs []int) (dataDates map[int]time.Time, err error) {
	action := "get dates by data IDs"

	dataDates = make(map[int]time.Time)

	for _, dataID := range dataIDs {
		row := rp.db.QueryRowContext(ctx, `
	
			SELECT
				created_at
			FROM
				data
			WHERE
				id = $1;
	
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

func (rp *postgresData) GetDatas(ctx context.Context, dataIDs []int) (datas []model.Data, err error) {
	action := "get datas by data IDs"

	datas = make([]model.Data, 0)

	for _, dataID := range dataIDs {
		row := rp.db.QueryRowContext(ctx, `
		
			SELECT
				group_id, data_type, title, created_at
			FROM
				data
			WHERE
				id = $1;
		
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
					descr, note
				FROM
					data_note
				WHERE
					data_id = $1;
			
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
					descr, login, pass
				FROM
					data_pass
				WHERE
					data_id = $1;
			
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
					descr, number, date, name, cvc2, pin, bank
				FROM
					data_card
				WHERE
					data_id = $1;
			
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
					descr, filename, filesize, filedate, file
				FROM
					data_file
				WHERE
					data_id = $1;
			
			`, dataID)

			if err = row2.Scan(&file.Desc, &file.Filename, &file.Filesize, &file.Filedate, &file.File); err != nil {
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

func (rp *postgresData) GetDataFile(ctx context.Context, dataID int) (file []byte, err error) {
	action := "get file by data ID"

	row := rp.db.QueryRowContext(ctx, `
		
		SELECT
			file
		FROM
			data_file
		WHERE
			data_id = $1;
		
	`, dataID)

	if err = row.Scan(&file); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return file, nil
}

func (rp *postgresData) AddDatas(ctx context.Context, datas []model.Data) (dataIDs map[int]int, err error) {
	action := "add datas"

	dataIDs = make(map[int]int)

	for _, data := range datas {
		row := rp.db.QueryRowContext(ctx, `
		
			INSERT INTO data
				(group_id, data_type, title, created_at)
			VALUES
				($1, $2, $3, $4)
			RETURNING
				id;
		
		`, data.GroupID, data.DataType, data.Title, data.CreatedAt)

		var dataID int
		if err = row.Scan(&dataID); err != nil {
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
			_, err = rp.db.ExecContext(ctx, `
			
				INSERT INTO data_note
					(data_id, descr, note)
				VALUES
					($1, $2, $3);
			
			`, dataID, data.Note.Desc, data.Note.Note)
			if err != nil {
				return nil, fmt.Errorf(
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
					(data_id, descr, login, pass)
				VALUES
					($1, $2, $3, $4);
			
			`, dataID, data.Pass.Desc, data.Pass.Login, data.Pass.Pass)
			if err != nil {
				return nil, fmt.Errorf(
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
					(data_id, descr, number, date, name, cvc2, pin, bank)
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8);
			
			`, dataID, data.Card.Desc, data.Card.Number, data.Card.Date, data.Card.Name, data.Card.CVC2, data.Card.PIN, data.Card.Bank)
			if err != nil {
				return nil, fmt.Errorf(
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
					(data_id, descr, filename, filesize, filedate)
				VALUES
					($1, $2, $3, $4, $5);
			
			`, dataID, data.File.Desc, data.File.Filename, data.File.Filesize, data.File.Filedate)
			if err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		}
		dataIDs[data.IDOnClient] = dataID
	}
	return dataIDs, nil
}

func (rp *postgresData) SaveDataFile(ctx context.Context, dataID int, file []byte) (err error) {
	action := "save file"

	_, err = rp.db.ExecContext(ctx, `
		
		UPDATE
			data_file
		SET
			file = $1
		WHERE
			data_id = $2;
		
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
