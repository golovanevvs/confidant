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

	if len(groupIDs) == 0 {
		return mapGroupIDEmails, nil
	}

	query, args, err := sqlx.In(`
	
		SELECT
			email, group_id
		FROM
			email_in_groups
		WHERE
			group_id IN (?)
		ORDER BY
			group_id;
	
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

	for rows.Next() {
		var email string
		var groupID int
		if err = rows.Scan(&email, &groupID); err != nil {
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

func (rp *postgresGroups) AddGroups(ctx context.Context, groups []model.Group) (mapGroupIDOnClientGroupID map[int]int, err error) {
	action := "add groups"

	mapGroupIDOnClientGroupID = make(map[int]int)

	if len(groups) == 0 {
		return mapGroupIDOnClientGroupID, nil
	}

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

	groupQuery := `
	
		INSERT INTO groups
			(id_on_client, title, account_id)
		VALUES
			(:id_on_client, :title, :account_id)
		RETURNING
			id, id_on_client;
	
	`

	var groupResults []struct {
		ID         int `db:"id"`
		IDOnClient int `db:"id_on_client"`
	}

	groupArgs := make([]map[string]interface{}, len(groups))
	for i, group := range groups {
		groupArgs[i] = map[string]interface{}{
			"id_on_client": group.IDOnClient,
			"title":        group.Title,
			"account_id":   group.AccountID,
		}
	}

	stmt, err := tx.PrepareNamedContext(ctx, groupQuery)
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

	for _, arg := range groupArgs {
		var res struct {
			ID         int `db:"id"`
			IDOnClient int `db:"id_on_client"`
		}
		if err = stmt.GetContext(ctx, &res, arg); err != nil {
			return nil, fmt.Errorf(
				"%s: %s: %w: %w",
				customerrors.DBErr,
				action,
				customerrors.ErrDBInternalError500,
				err,
			)
		}
		groupResults = append(groupResults, res)
	}

	if err = rp.batchInsertEmails(ctx, tx, groupResults, groups); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w: %w",
			customerrors.DBErr,
			action,
			customerrors.ErrDBInternalError500,
			err,
		)
	}

	for _, res := range groupResults {
		mapGroupIDOnClientGroupID[res.IDOnClient] = res.ID
	}

	return mapGroupIDOnClientGroupID, nil

	// for _, group := range groups {
	// 	row := rp.db.QueryRowContext(ctx, `

	// 	INSERT INTO groups
	// 		(title, account_id)
	// 	VALUES
	// 		($1, $2)
	// 	RETURNING id;

	// `, group.Title, group.AccountID)

	// 	var groupID int
	// 	if err = row.Scan(&groupID); err != nil {
	// 		return nil, fmt.Errorf(
	// 			"%s: %s: %w: %w",
	// 			customerrors.DBErr,
	// 			action,
	// 			customerrors.ErrDBInternalError500,
	// 			err,
	// 		)
	// 	}

	// 	for _, email := range group.Emails {
	// 		_, err := rp.db.ExecContext(ctx, `

	// 		INSERT INTO email_in_groups
	// 			(email, group_id)
	// 		VALUES
	// 			($1, $2);

	// 	`, email, groupID)

	// 		if err != nil {
	// 			return nil, fmt.Errorf(
	// 				"%s: %s: %w: %w",
	// 				customerrors.DBErr,
	// 				action,
	// 				customerrors.ErrDBInternalError500,
	// 				err,
	// 			)
	// 		}
	// 	}

	// 	groupIDs[group.IDOnClient] = groupID
	// }

	// return groupIDs, nil
}

func (rp *postgresGroups) batchInsertEmails(ctx context.Context, tx *sqlx.Tx, groupResults []struct {
	ID         int `db:"id"`
	IDOnClient int `db:"id_on_client"`
}, groups []model.Group) (err error) {
	action := "batch insert emails"
	var emailArgs []map[string]interface{}
	for _, res := range groupResults {
		for _, group := range groups {
			if group.IDOnClient == res.IDOnClient {
				for _, email := range group.Emails {
					emailArgs = append(emailArgs, map[string]interface{}{
						"group_id": res.ID,
						"email":    email,
					})
				}
				break
			}
		}
	}
	if len(emailArgs) == 0 {
		return nil
	}

	_, err = tx.NamedExecContext(ctx, `
	
		INSERT INTO email_in_groups
			(group_id, email)
		VALUES
			(:group_id, :email);
	
	`, emailArgs)
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
