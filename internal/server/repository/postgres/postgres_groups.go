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

func (rp *postgresGroups) GetGroupIDs(ctx context.Context, accountID int) (groupIDs []int, err error) {
	action := "get group IDs"

	groupIDs = make([]int, 0)

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			email_in_groups.group_id
		FROM
			email_in_groups
		JOIN
			account
		ON
			email_in_groups.email = account.email
		WHERE
			account.id = $1;
	
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
		groupIDs = append(groupIDs, groupID)
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

func (rp *postgresGroups) GetGroups(ctx context.Context, accountID int, groupIDs []int) (groups []model.Group, err error) {
	action := "get groups"

	if len(groupIDs) == 0 {
		return []model.Group{}, nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			g.id, g.title, g.account_id, e.email
		FROM
			groups g
		LEFT JOIN
			email_in_groups e
		ON
			g.id = e.group_id
		WHERE
			g.id IN (?)
		ORDER BY
			g.id;
	
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

	rows, err := rp.db.QueryContext(ctx, rp.db.Rebind(query), args...)
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

	mapGroupIDGroup := make(map[int]*model.Group)
	for rows.Next() {
		var (
			groupID   int
			title     string
			accountID int
			email     sql.NullString
		)
		if err = rows.Scan(&groupID, &title, &accountID, &email); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		group, inMap := mapGroupIDGroup[groupID]
		if !inMap {
			group = &model.Group{
				ID:        groupID,
				Title:     title,
				AccountID: accountID,
				Emails:    []string{},
			}
			mapGroupIDGroup[groupID] = group
		}
		if email.Valid {
			group.Emails = append(group.Emails, email.String)
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

	groups = make([]model.Group, 0, len(mapGroupIDGroup))
	for _, group := range mapGroupIDGroup {
		groups = append(groups, *group)
	}

	return groups, nil
}

func (rp *postgresGroups) GetEmails(ctx context.Context, groupIDs []int) (mapGroupIDEmails map[int][]string, err error) {
	action := "get map groupID-emails by group IDs"

	mapGroupIDEmails = make(map[int][]string, 0)

	for _, groupID := range groupIDs {
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

		mapGroupIDEmails[groupID] = emails
	}

	return mapGroupIDEmails, nil
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
			return nil, fmt.Errorf(
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
				return nil, fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrDBInternalError500,
					err,
				)
			}
		}

		groupIDs[group.IDOnClient] = groupID
	}

	return groupIDs, nil
}

func (rp *postgresGroups) AddEmails(ctx context.Context, mapGroupIDEmails map[int][]string) (err error) {
	action := "add e-mails"

	for groupID, emails := range mapGroupIDEmails {
		for _, email := range emails {
			_, err = rp.db.ExecContext(ctx, `
	
				INSERT INTO email_in_groups
					(group_id, email)
				VALUES
					($1, $2);
	
	`, groupID, email)

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
