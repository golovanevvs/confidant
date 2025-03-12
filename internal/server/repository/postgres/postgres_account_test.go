package postgres

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAccount(t *testing.T) {
	//! preparatory operations
	ctx := context.Background()

	//! using a mock DB
	// creating a mock db
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %s", err.Error())
	}
	defer db.Close()

	//creating an sqlx wrapper
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	//initializing the accountPostgres
	rp := NewPostgresAccount(sqlxDB)

	// configuring expected values
	mock.ExpectQuery(`
	
		^INSERT INTO account 
			\(email, password_hash\)
		VALUES
			\(\$1, \$2\)
	 	RETURNING id;$
	 
	 `).WithArgs("test@test.com", "testPasswordHash").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	//! runnig the tests
	accountID, err := rp.SaveAccount(ctx,
		model.Account{
			Email:        "test@test.com",
			PasswordHash: "testPasswordHash",
		},
	)
	require.NoError(t, err)
	assert.Equal(t, 1, accountID)
}
