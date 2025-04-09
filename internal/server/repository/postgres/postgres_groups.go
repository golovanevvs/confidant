package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/jmoiron/sqlx"
)

type postgresGroups struct {
	db *sqlx.DB
}

func NewPostgresGroups(db *sqlx.DB) *postgresGroups {
	return &postgresGroups{
		db: db,
	}
}

func (rp *postgresGroups) GetGroupIDs(ctx context.Context, accountID int) (groupIDs map[int]struct{}, err error) {
	action := "get group IDs"

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
			email = $1
		GROUP BY
			group_id;
	
	`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return groupIDs, nil
		} else {
			return groupIDs, fmt.Errorf(
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
		return groupIDs, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return groupIDs, nil
}

func (rp *postgresGroups) GetGroups(ctx context.Context, accountID int, groupIDs []int) (groups []model.Group, err error) {
	action := "get groups"

	for _, groupID := range groupIDs {

		row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			title, account_id
		FROM
			groups
		WHERE
			id = $1;
	
	`, groupID)

		var group model.Group
		if err = row.Scan(&group.Title, &group.AccountID); err != nil {
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
					email
				FROM
					email_in_groups
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

		emails := make([]string, 0)

		for rows.Next() {
			var email string
			if err = rows.Scan(&email); err != nil {
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
			emails = append(emails, email)
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

		group.ID = groupID
		group.Emails = emails

		groups = append(groups, group)
	}

	return groups, nil
}

func (rp *postgresGroups) AddGroups(ctx context.Context, groups []model.Group) (groupIDs map[int]int, err error) {
	action := "add groups"

	groupIDs = make(map[int]int)

	for _, group := range groups {
		row := rp.db.QueryRowContext(ctx, `
	
		INSERT INTO groups
			(title, account_id)
		VALUES
			($1, $2)
		RETURNING id;
	
	`, group.Title, group.AccountID)

		var groupID int
		if err = row.Scan(&groupID); err != nil {
			return groupIDs, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		for _, email := range group.Emails {
			_, err := rp.db.ExecContext(ctx, `
		
			INSERT INTO email_in_groups
				(email, group_id)
			VALUES
				($1, $2);
		
		`, email, groupID)

			if err != nil {
				return groupIDs, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		}

		groupIDs[group.ID] = groupID
	}

	return groupIDs, nil
}
