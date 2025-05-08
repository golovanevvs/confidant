package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	*postgresManage
	*postgresAccount
	*postgresGroups
	*postgresData
}

func New(databaseURI string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	rpgp := NewPostgresGroups(db)

	return &Postgres{
		NewPostgresManage(db),
		NewPostgresAccount(db),
		rpgp,
		NewPostgresData(db, rpgp),
	}, nil
}
