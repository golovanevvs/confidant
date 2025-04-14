package postgres

import (
	"fmt"

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
	db, err := sqlx.Open("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("DB ping error: %s", err.Error())
	}

	rpgp := NewPostgresGroups(db)

	return &Postgres{
		NewPostgresManage(db),
		NewPostgresAccount(db),
		rpgp,
		NewPostgresData(db, rpgp),
	}, nil
}
