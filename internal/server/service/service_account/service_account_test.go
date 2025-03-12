package service_account

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golovanevvs/confidant/internal/customerrors"
	"github.com/golovanevvs/confidant/internal/server/model"
	"github.com/golovanevvs/confidant/internal/server/service/repository_mock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount(t *testing.T) {
	//! preparatory operations
	ctx := context.Background()

	//! using a repository mock
	// creating a gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creating a repository mock-object
	rpMock := repository_mock.NewMockIRepository(ctrl)

	//initializing the service
	sv := New(rpMock)

	//! setting input and output parameters
	// input parameters
	type input struct {
		account model.Account
	}

	// expected output values
	type expected struct {
		accountID int
		err       error
	}

	tests := []struct {
		name      string   // test name
		input     input    // input parameters
		setupMock func()   // setup mock
		expected  expected // expected output values
	}{
		{
			name: "successful creating",
			input: input{
				account: model.Account{
					Email:    "test@test.ru",
					Password: "Ul34.fb2",
				},
			},
			setupMock: func() {
				rpMock.EXPECT().
					SaveAccount(
						gomock.Any(),
						gomock.Any(),
					).Return(1, nil)
			},
			expected: expected{
				accountID: 1,
				err:       nil,
			},
		},
		{
			name: "error creating",
			input: input{
				account: model.Account{
					Email:    "test@test.ru",
					Password: "Ul34.fb2",
				},
			},
			setupMock: func() {
				rpMock.EXPECT().
					SaveAccount(
						gomock.Any(),
						gomock.Any(),
					).Return(-1, fmt.Errorf("%s: %s: %w", customerrors.DBErr, "save account", customerrors.ErrDBBusyEmail409))
			},
			expected: expected{
				accountID: -1,
				err:       fmt.Errorf("%s", customerrors.DBErr),
			},
		},
	}

	//! runnig the tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setuping mock
			test.setupMock()

			accountID, err := sv.CreateAccount(ctx, test.input.account)

			if test.expected.err != nil {
				assert.ErrorContains(t, err, test.expected.err.Error())
			}

			assert.Equal(t, test.expected.accountID, accountID)
		})
	}

}
