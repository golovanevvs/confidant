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
		for _, groupID := range groupIDs {
			rows, err := rp.db.QueryContext(ctx, `
			
				SELECT
					id
				FROM
					data
				WHERE
					group_id = $1;
	
			`, groupID)

			if err != nil && err != sql.ErrNoRows {
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
