package repository

import "sync"

type (
	CartRepository interface {
		GetCart(clientId int)
	}
	CartRepositoryImpl struct{}
)

var (
	once     sync.Once
	dbClient DBClient
)
