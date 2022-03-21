package repository

import (
	"github.com/emiliocc5/online-store-api/internal/models"
	"gorm.io/gorm"
)

type (
	ClientRepository interface {
		IsClientInDataBase(clientId int) bool
	}
	PgClientRepository struct {
		DbClient IClientRepositoryDbClient
	}
	IClientRepositoryDbClient interface {
		Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	}
)

func (cr *PgClientRepository) IsClientInDataBase(clientId int) bool {
	client := models.Client{}
	findClientResult := cr.DbClient.Find(&client, "id = ?", clientId)
	return findClientResult.Error == nil
}
