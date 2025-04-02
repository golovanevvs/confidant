package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type postgresGroups struct {
	db *sqlx.DB
}

func NewPostgresGroups(db *sqlx.DB) *postgresManage {
	return &postgresManage{
		db: db,
	}
}

func (rp *postgresGroups) GetGroupIDs(ctx context.Context, accountID int) (groupIDs map[int]struct{}, err error) {
	action := "get groups"

	groupIDs = make(map[int]struct{})

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			email
		FROM
			account
		WHERE
			account_id = $1;
	
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
		groupIDs[groupID] = struct{}{}
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

	return groupIDs, nil
}
