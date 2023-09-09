package store

import (
	"math/big"
	"sync"

	"github.com/jayendramadaram/port-wardens/server"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	server.Store
}

type store struct {
	mu    *sync.RWMutex
	cache map[string]*big.Float

	*mongo.Database
}

func NewStore(db *mongo.Database) Store {
	return &store{new(sync.RWMutex), make(map[string]*big.Float), db}
}

func (s *store) HealthCheck() error {
	return nil
}
