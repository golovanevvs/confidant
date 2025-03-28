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
	
		INSERT INTO email_in_groups
			(email, group_id)
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

func (rp *sqliteGroups) GetGroups(ctx context.Context, email string) (groups []model.Group, err error) {
	action := "get groups"

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			groups.id, groups.id_on_server, groups.title
		FROM
			groups
		INNER JOIN
			email_in_groups
		ON
			groups.id = email_in_groups.group_id
		WHERE
			email = ?;

	`, email)
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

func (rp *sqliteData) GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error) {
	action := "get group ID"

	row := rp.db.QueryRowContext(ctx, `
	
		SELECT
			groups.id
		FROM
			groups
		INNER JOIN
			email_in_groups
		ON
			groups.id = email_in_groups.group_id
		WHERE
			email_in_groups.email = ? AND groups.title = ?;
	
	`, email, titleGroup)

	if err = row.Scan(&groupID); err != nil {
		return -1, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return groupID, nil
}

func (rp *sqliteGroups) getEmails(ctx context.Context, groupID int) (emails []string, err error) {
	action := "get emails"

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			email
		FROM
			email_in_groups
		WHERE
			group_id = ?;

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
			"%s 7: %w: %w",
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return emails, nil
}

func (rp *sqliteGroups) AddEmail(ctx context.Context, groupID int, email string) (err error) {
	action := "add e-mail"

	_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO email_in_groups
			(email, group_id)
		VALUES
			(?,?);
	
	`, email, groupID)
	if err != nil {
		return fmt.Errorf(
			"%s: %w: %w",
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	return nil
}
