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

	if len(groupIDs) == 0 {
		return nil, nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			id
		FROM
			data
		WHERE
			group_id IN (?);
		
		`, groupIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	query = rp.db.Rebind(query)

	rows, err := rp.db.QueryContext(ctx, query, args...)
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

	dataIDs = make([]int, 0)
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

	return dataIDs, nil
}

func (rp *postgresData) GetDataDates(ctx context.Context, dataIDs []int) (mapDataIDDate map[int]time.Time, err error) {
	action := "get dates by data IDs"

	query, args, err := sqlx.In(`
	
		SELECT
			id, created_at
		FROM
			data
		WHERE
			id IN (?);
	
	`, dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	query = rp.db.Rebind(query)

	rows, err := rp.db.QueryContext(ctx, query, args...)
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

	mapDataIDDate = make(map[int]time.Time)

	var dataID int
	var date time.Time
	for rows.Next() {
		if err = rows.Scan(&dataID, &date); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
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

	mapDataIDDate[dataID] = date

	return mapDataIDDate, nil
}

func (rp *postgresData) GetDatas(ctx context.Context, dataIDs []int) (datas []model.Data, err error) {
	action := "get datas by data IDs"

	if len(dataIDs) == 0 {
		return []model.Data{}, nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			id, group_id, data_type, title, created_at
		FROM
			data
		WHERE
			id IN (?);
	
	`, dataIDs)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	err = rp.db.SelectContext(ctx, &datas, rp.db.Rebind(query), args...)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	if err = rp.loadNoteData(ctx, &datas); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	if err = rp.loadPassData(ctx, &datas); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	if err = rp.loadCardData(ctx, &datas); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	if err = rp.loadFileData(ctx, &datas); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return datas, nil
}

func (rp *postgresData) loadNoteData(ctx context.Context, datas *[]model.Data) error {
	action := "load note datas"

	var noteDataIDs []int

	for _, data := range *datas {
		if data.DataType == "note" {
			noteDataIDs = append(noteDataIDs, data.ID)
		}
	}

	if len(noteDataIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			data_id, descr, note
		FROM
			data_note
		WHERE
			data_id IN (?);
	
	`, noteDataIDs)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	notesEnc := []model.NoteEnc{}

	if err = rp.db.SelectContext(ctx, &notesEnc, rp.db.Rebind(query), args...); err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	mapDataIDNoteEnc := make(map[int]model.NoteEnc)
	for _, note := range notesEnc {
		mapDataIDNoteEnc[note.DataID] = model.NoteEnc{
			Desc: note.Desc,
			Note: note.Note,
		}
	}

	for i := range *datas {
		data := &(*datas)[i]
		if data.DataType == "note" {
			if note, inMap := mapDataIDNoteEnc[data.ID]; inMap {
				data.Note = note
			} else {
				return fmt.Errorf("%s: %w", action, err)
			}
		}
	}

	return nil
}

func (rp *postgresData) loadPassData(ctx context.Context, datas *[]model.Data) error {
	action := "load pass datas"

	var passDataIDs []int

	for _, data := range *datas {
		if data.DataType == "pass" {
			passDataIDs = append(passDataIDs, data.ID)
		}
	}

	if len(passDataIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			data_id, descr, login, pass
		FROM
			data_pass
		WHERE
			data_id IN (?);
	
	`, passDataIDs)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	passesEnc := []model.PassEnc{}

	err = rp.db.SelectContext(ctx, &passesEnc, rp.db.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	mapDataIDPassEnc := make(map[int]model.PassEnc)
	for _, pass := range passesEnc {
		mapDataIDPassEnc[pass.DataID] = model.PassEnc{
			Desc:  pass.Desc,
			Login: pass.Login,
			Pass:  pass.Pass,
		}
	}

	for i := range *datas {
		data := &(*datas)[i]
		if data.DataType == "pass" {
			if pass, inMap := mapDataIDPassEnc[data.ID]; inMap {
				data.Pass = pass
			} else {
				return fmt.Errorf("%s: %w", action, err)
			}
		}
	}

	return nil
}

func (rp *postgresData) loadCardData(ctx context.Context, datas *[]model.Data) error {
	action := "load card datas"

	var cardDataIDs []int

	for _, data := range *datas {
		if data.DataType == "card" {
			cardDataIDs = append(cardDataIDs, data.ID)
		}
	}

	if len(cardDataIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			data_id, descr, number, date, name, cvc2, pin, bank
		FROM
			data_card
		WHERE
			data_id IN (?);
	
	`, cardDataIDs)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	cardsEnc := []model.CardEnc{}

	err = rp.db.SelectContext(ctx, &cardsEnc, rp.db.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	mapDataIDCardEnc := make(map[int]model.CardEnc)
	for _, card := range cardsEnc {
		mapDataIDCardEnc[card.DataID] = model.CardEnc{
			Desc:   card.Desc,
			Number: card.Number,
			Date:   card.Date,
			Name:   card.Name,
			CVC2:   card.CVC2,
			PIN:    card.PIN,
			Bank:   card.Bank,
		}
	}

	for i := range *datas {
		data := &(*datas)[i]
		if data.DataType == "card" {
			if card, inMap := mapDataIDCardEnc[data.ID]; inMap {
				data.Card = card
			} else {
				return fmt.Errorf("%s: %w", action, err)
			}
		}
	}

	return nil
}

func (rp *postgresData) loadFileData(ctx context.Context, datas *[]model.Data) error {
	action := "load file datas"

	var fileDataIDs []int

	for _, data := range *datas {
		if data.DataType == "file" {
			fileDataIDs = append(fileDataIDs, data.ID)
		}
	}

	if len(fileDataIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			data_id, descr, filename, filesize, filedate
		FROM
			data_file
		WHERE
			data_id IN (?);
	
	`, fileDataIDs)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	filesEnc := []model.FileEnc{}

	err = rp.db.SelectContext(ctx, &filesEnc, rp.db.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("%s: %w", action, err)
	}

	mapDataIDFileEnc := make(map[int]model.FileEnc)
	for _, file := range filesEnc {
		mapDataIDFileEnc[file.DataID] = model.FileEnc{
			Desc:     file.Desc,
			Filename: file.Filename,
			Filesize: file.Filesize,
			Filedate: file.Filedate,
		}
	}

	for i := range *datas {
		data := &(*datas)[i]
		if data.DataType == "file" {
			if file, inMap := mapDataIDFileEnc[data.ID]; inMap {
				data.File = file
			} else {
				return fmt.Errorf("%s: %w", action, err)
			}
		}
	}

	return nil
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

func (rp *postgresData) AddDatas(ctx context.Context, datas []model.Data) (mapDataIDOnClientDataID map[int]int, err error) {
	action := "add datas"

	mapDataIDOnClientDataID = make(map[int]int)

	tx, err := rp.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	dataQuery := `
	
		INSERT INTO data
			(id_on_client, group_id, data_type, title, created_at)
		VALUES
			(:id_on_client, :group_id, :data_type, :title, :created_at)
		RETURNING
			id, id_on_client;
	
	`

	dataArgs := make([]map[string]interface{}, len(datas))
	for i, data := range datas {
		dataArgs[i] = map[string]interface{}{
			"id_on_client": data.IDOnClient,
			"group_id":     data.GroupID,
			"data_type":    data.DataType,
			"title":        data.Title,
			"created_at":   data.CreatedAt,
		}
	}

	stmt, err := tx.PrepareNamedContext(ctx, dataQuery)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	defer stmt.Close()

	var dataResults []model.DataResults
	for _, arg := range dataArgs {
		var dataResult model.DataResults
		if err = stmt.GetContext(ctx, &dataResult, arg); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		dataResults = append(dataResults, dataResult)
		mapDataIDOnClientDataID[dataResult.IDOnClient] = dataResult.ID
	}

	if err = rp.batchInsertNotes(ctx, tx, datas, dataResults); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	if err = rp.batchInsertPasses(ctx, tx, datas, dataResults); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	if err = rp.batchInsertCards(ctx, tx, datas, dataResults); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	if err = rp.batchInsertFiles(ctx, tx, datas, dataResults); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return mapDataIDOnClientDataID, nil
}

func (rp *postgresData) batchInsertNotes(ctx context.Context, tx *sqlx.Tx, datas []model.Data, dataResults []model.DataResults) (err error) {
	action := "batch insert notes"

	var notes []map[string]interface{}

	for _, dataResult := range dataResults {
		for _, data := range datas {
			if data.IDOnClient == dataResult.IDOnClient && data.DataType == "note" {
				notes = append(notes, map[string]interface{}{
					"data_id": dataResult.ID,
					"descr":   data.Note.Desc,
					"note":    data.Note.Note,
				})
				break
			}
		}
	}

	if len(notes) > 0 {
		_, err = tx.NamedExecContext(ctx, `
		
			INSERT INTO data_note
				(data_id, descr, note)
			VALUES
				(:data_id, :descr, :note);
		
		`, notes)
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

	return nil
}

func (rp *postgresData) batchInsertPasses(ctx context.Context, tx *sqlx.Tx, datas []model.Data, dataResults []model.DataResults) (err error) {
	action := "batch insert passes"

	var passes []map[string]interface{}

	for _, dataResult := range dataResults {
		for _, data := range datas {
			if data.IDOnClient == dataResult.IDOnClient && data.DataType == "pass" {
				passes = append(passes, map[string]interface{}{
					"data_id": dataResult.ID,
					"descr":   data.Pass.Desc,
					"login":   data.Pass.Login,
					"pass":    data.Pass.Pass,
				})
				break
			}
		}
	}

	if len(passes) > 0 {
		_, err = tx.NamedExecContext(ctx, `
		
			INSERT INTO data_pass
				(data_id, descr, login, pass)
			VALUES
				(:data_id, :descr, :login, :pass);
		
		`, passes)
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
	return nil
}

func (rp *postgresData) batchInsertCards(ctx context.Context, tx *sqlx.Tx, datas []model.Data, dataResults []model.DataResults) (err error) {
	action := "batch insert cards"

	var cards []map[string]interface{}

	for _, dataResult := range dataResults {
		for _, data := range datas {
			if data.IDOnClient == dataResult.IDOnClient && data.DataType == "card" {
				cards = append(cards, map[string]interface{}{
					"data_id": dataResult.ID,
					"descr":   data.Card.Desc,
					"number":  data.Card.Number,
					"date":    data.Card.Date,
					"name":    data.Card.Name,
					"cvc2":    data.Card.CVC2,
					"pin":     data.Card.PIN,
					"bank":    data.Card.Bank,
				})
				break
			}
		}
	}

	if len(cards) > 0 {
		_, err = tx.NamedExecContext(ctx, `
		
			INSERT INTO data_card
				(data_id, descr, number, date, name, cvc2, pin, bank)
			VALUES
				(:data_id, :descr, :number, :date, :name, :cvc2, :pin, :bank);
		
		`, cards)
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
	return nil
}

func (rp *postgresData) batchInsertFiles(ctx context.Context, tx *sqlx.Tx, datas []model.Data, dataResults []model.DataResults) (err error) {
	action := "batch insert files"

	var files []map[string]interface{}

	for _, dataResult := range dataResults {
		for _, data := range datas {
			if data.IDOnClient == dataResult.IDOnClient && data.DataType == "file" {
				files = append(files, map[string]interface{}{
					"data_id":  dataResult.ID,
					"descr":    data.File.Desc,
					"filename": data.File.Filename,
					"filesize": data.File.Filesize,
					"filedate": data.File.Filedate,
				})
				break
			}
		}
	}

	if len(files) > 0 {
		_, err = tx.NamedExecContext(ctx, `
		
			INSERT INTO data_file
				(data_id, descr, filename, filesize, filedate)
			VALUES
				(:data_id, :descr, :filename, :filesize, :filedate);
		
		`, files)
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
	return nil
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
