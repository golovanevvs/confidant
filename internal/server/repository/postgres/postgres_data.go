package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type postgresData struct {
	db *sqlx.DB
}

func NewPostgresData(db *sqlx.DB) *postgresData {
	return &postgresData{
		db: db,
	}
}

func (rp *postgresGroups) GetDataIDs(ctx context.Context, accountID int) (dataIDs []int, err error) {
	action := "get data IDs"

	// getting group IDs
	var groupIDs []int
	groupIDs, err = rp.GetGroupIDs(ctx, accountID)
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
		//! ----------------- СТОП --------------------

		row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			email
		FROM
			account
		WHERE
			id = $1;
	
	`, accountID)

		var email string
		if err = row.Scan(&email); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			group_id
		FROM
			email_in_groups
		WHERE
			email = $1;
	
	`, email)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			} else {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		}
		defer rows.Close()

		var groupIDs2 []int
		for rows.Next() {
			var groupID int
			if err = rows.Scan(&groupID); err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
			groupIDs2 = append(groupIDs2, groupID)
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
	} else {
		return nil, nil
	}
}
