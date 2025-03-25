package db_sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/golovanevvs/confidant/internal/client/model"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/jmoiron/sqlx"
)

type sqliteGroups struct {
	db *sqlx.DB
}

func NewSQLiteGroups(db *sqlx.DB) *sqliteGroups {
	return &sqliteGroups{
		db: db,
	}
}

func (rp *sqliteGroups) AddGroup(ctx context.Context, account *model.Account, title string) (err error) {
	action := "add group"

	row := rp.db.QueryRowContext(ctx, `

		INSERT INTO groups
			(title, account_id)
		VALUES
			(?, ?)
		RETURNING id;

	`, title, account.ID)
	var groupID int
	err = row.Scan(&groupID)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrAddGroup,
				err,
			)
		} else {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
	}

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO emails_in_groups
			(email, groups_id)
		VALUES
			(?, ?);

	`, account.Email, groupID)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddEmailInGroup,
			err,
		)
	}

	return nil
}

func (rp *sqliteGroups) GetGroups(ctx context.Context, accountID int) (groups []model.Group, err error) {
	action := "get groups"

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			id, id_on_server, title
		FROM
			groups
		WHERE
			account_id = ?;

	`, accountID)
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
		var group model.Group
		if err = rows.Scan(&group.ID, &group.IDOnServer, &group.Title); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		emails, err := rp.getEmails(ctx, group.ID)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		group.Emails = emails

		groups = append(groups, group)
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

	return groups, nil
}

func (rp *sqliteGroups) getEmails(ctx context.Context, groupID int) (emails []string, err error) {
	action := "get emails"

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			email
		FROM
			emails_in_groups
		WHERE
			groups_id = ?;

	`, groupID)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %w: %w",
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var email string
		if err = rows.Scan(&email); err != nil {
			return nil, fmt.Errorf(
				"%s: %w: %w",
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		emails = append(emails, email)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"%s: %w: %w",
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return emails, nil
}
