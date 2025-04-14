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
			ID:   dataID,
			Note: note,
			Pass: pass,
			Card: card,
			File: file,
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
					data_note
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
