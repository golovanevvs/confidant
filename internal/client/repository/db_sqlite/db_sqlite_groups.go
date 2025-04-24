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

	tx, err := rp.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddGroup,
			err,
		)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	row := tx.QueryRowContext(ctx, `

		INSERT INTO groups
			(title, account_id)
		VALUES
			(?, ?)
		RETURNING id;

	`, title, account.ID)

	var groupID int
	err = row.Scan(&groupID)
	if err != nil {
		_ = tx.Rollback()
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

	_, err = tx.ExecContext(ctx, `
	
		INSERT INTO email_in_groups
			(email, group_id)
		VALUES
			(?, ?);

	`, account.Email, groupID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddEmailInGroup,
			err,
		)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrAddGroup,
			err,
		)
	}

	return nil
}

func (rp *sqliteGroups) GetGroups(ctx context.Context, email string) (groups []model.Group, err error) {
	action := "get groups"

	groupIDs := make([]int, 0)

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
			email IN (?);

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

		groupIDs = append(groupIDs, group.ID)
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

	emailQuery, emailArgs, err := sqlx.In(`
		SELECT
			email, group_id
		FROM
			email_in_groups
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

	emailRows, err := rp.db.QueryContext(ctx, emailQuery, emailArgs...)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}
	defer emailRows.Close()

	mapGroupIDEmails := make(map[int][]string, 0)
	for emailRows.Next() {
		var email string
		var groupID int
		if err = emailRows.Scan(&email, &groupID); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		mapGroupIDEmails[groupID] = append(mapGroupIDEmails[groupID], email)
	}

	for i := range groups {
		groups[i].Emails = mapGroupIDEmails[groups[i].ID]
	}

	return groups, nil
}

func (rp *sqliteGroups) GetEmails(ctx context.Context, groupIDs []int) (mapGroupIDEmails map[int][]string, err error) {
	action := "get map groupID-emails by group IDs"

	mapGroupIDEmails = make(map[int][]string, 0)

	query, args, err := sqlx.In(`
			
		SELECT
			email_in_groups.email, groups.id_on_server
		FROM
			email_in_groups
		INNER JOIN
			groups
		ON
			email_in_groups.group_id = groups.id
		WHERE
			groups.id_on_server IN (?);
			
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

	for rows.Next() {
		var email string
		var groupOnServerID int
		if err = rows.Scan(&email, &groupOnServerID); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		mapGroupIDEmails[groupOnServerID] = append(mapGroupIDEmails[groupOnServerID], email)

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

	return mapGroupIDEmails, nil
}

func (rp *sqliteGroups) GetGroupsByIDs(ctx context.Context, groupIDs []int) (groups []model.Group, err error) {
	action := "get groups by IDs"

	query, args, err := sqlx.In(`
	
		SELECT
			id, title, account_id
		FROM
			groups
		WHERE
			id = ?;

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

	for rows.Next() {
		var group model.Group
		if err = rows.Scan(&group.ID, &group.Title, &group.AccountID); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
	}

	emails, err := rp.getEmails(ctx, groupID)
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
	group.ID = groupID

	groups = append(groups, group)

	return groups, nil
}

func (rp *sqliteGroups) GetGroupID(ctx context.Context, email string, titleGroup string) (groupID int, err error) {
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
			"%s: %s: %w: %w, e-mail: %s, titleGroup: %s",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
			email,
			titleGroup,
		)
	}

	return groupID, nil
}

// func (rp *sqliteGroups) getEmails(ctx context.Context, groupID int) (emails []string, err error) {
// 	action := "get emails"

// 	emails = make([]string, 0)

// 	rows, err := rp.db.QueryContext(ctx, `

// 		SELECT
// 			email
// 		FROM
// 			email_in_groups
// 		WHERE
// 			group_id = ?;

// 	`, groupID)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"%s: %w: %w",
// 			action,
// 			customerrors.ErrDBInternalError500,
// 			err,
// 		)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var email string
// 		if err = rows.Scan(&email); err != nil {
// 			return nil, fmt.Errorf(
// 				"%s: %w: %w",
// 				action,
// 				customerrors.ErrDBInternalError500,
// 				err,
// 			)
// 		}

// 		emails = append(emails, email)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf(
// 			"%s 7: %w: %w",
// 			action,
// 			customerrors.ErrDBInternalError500,
// 			err,
// 		)
// 	}

// 	return emails, nil
// }

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

func (rp *sqliteGroups) GetGroupIDs(ctx context.Context, email string) (groupServerIDs []int, groupNoServerIDs []int, err error) {
	action := "get group IDs"

	groupServerIDs = make([]int, 0)
	groupNoServerIDs = make([]int, 0)

	rows, err := rp.db.QueryContext(ctx, `
	
		SELECT
			groups.id, groups.id_on_server
		FROM
			groups
		INNER JOIN
			email_in_groups
		ON
			groups.id = email_in_groups.group_id
		WHERE
			email_in_groups.email = ?;
	
	`, email)
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
		var groupID, groupServerID int
		if err = rows.Scan(&groupID, &groupServerID); err != nil {
			return nil, nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}

		if groupServerID == -1 {
			groupNoServerIDs = append(groupNoServerIDs, groupID)
		} else {
			groupServerIDs = append(groupServerIDs, groupServerID)
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

	return groupServerIDs, groupNoServerIDs, nil
}

func (rp *sqliteGroups) AddGroupBySync(ctx context.Context, group model.Group) (err error) {
	action := "add group by sync"

	row := rp.db.QueryRowContext(ctx, `

		INSERT INTO groups
			(id_on_server, title, account_id)
		VALUES
			(?, ?, ?)
		RETURNING id;

	`, group.IDOnServer, group.Title, group.AccountID)

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

	for _, email := range group.Emails {

		_, err = rp.db.ExecContext(ctx, `
	
		INSERT INTO email_in_groups
			(email, group_id)
		VALUES
			(?, ?);

	`, email, groupID)
		if err != nil {
			return fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrAddEmailInGroup,
				err,
			)
		}
	}

	return nil
}

func (rp *sqliteGroups) AddEmailsBySync(ctx context.Context, mapGroupIDEmails map[int][]string) (err error) {
	action := "add e-mails by sync"

	for groupIDOnServer, emails := range mapGroupIDEmails {
		for _, email := range emails {

			_, err = rp.db.ExecContext(ctx, `

				INSERT INTO email_in_groups
					(email, group_id)
				VALUES
					(?,
					(
						SELECT
							id
						FROM
							groups
						WHERE
							id_on_server = ?
					)
					);

				`, email, groupIDOnServer)

			if err != nil {
				return fmt.Errorf(
					"%s: %s: %w: %w",
					customerrors.DBErr,
					action,
					customerrors.ErrAddGroup,
					err,
				)
			}
		}
	}

	return nil
}

func (rp *sqliteGroups) UpdateGroupIDsOnServer(ctx context.Context, newGroupIDs map[int]int) (err error) {
	action := "update group IDsOnServer"

	for groupID, groupIDOnServer := range newGroupIDs {
		_, err = rp.db.ExecContext(ctx, `

		UPDATE
			groups
		SET
			id_on_server = ?
		WHERE
			id = ?;

	`, groupIDOnServer, groupID)
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
