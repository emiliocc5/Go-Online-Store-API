package repository

import (
	"errors"
	"github.com/emiliocc5/online-store-api/internal/models"
	"github.com/emiliocc5/online-store-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_GivenAValidClientId_ThenReturnTrue(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(getMockedClient(), nil)).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.Client)
		arg.Name = getMockedClient().Name
		arg.Id = getMockedClient().Id
	})

	repo := repository.PgClientRepository{DbClient: dbClientMock}

	resp := repo.IsClientInDataBase(aValidClientId)

	assert.True(t, resp)
}

func Test_GivenANotValidClientId_ThenReturnFalse(t *testing.T) {
	dbClientMock := &DbClientMock{}

	findClientByIdQuery := getMockedQuery("id = ?", aNotValidClientId)

	client := models.Client{}
	dbClientMock.On("Find", &client, findClientByIdQuery).Return(getMockedDbObject(nil, errors.New("client not found in db")))

	repo := repository.PgClientRepository{DbClient: dbClientMock}

	resp := repo.IsClientInDataBase(aNotValidClientId)

	assert.False(t, resp)
}
