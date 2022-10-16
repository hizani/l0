package cache

import (
	"wbintern/l0/service/database"
	"wbintern/l0/service/model"
)

// Cache for orders
type Cache map[string]model.OrderModel

// Initialize cache
func New() Cache {
	c := make(Cache)
	return c
}

// Restore cache from a database
func (c Cache) Restore(db *database.Database) error {
	orders, err := db.GetOrders()
	if err != nil {
		return err
	}
	for _, order := range orders {
		c.Add(order)
	}
	return err
}

// Add new order into cache
func (c Cache) Add(om model.OrderModel) {
	c[om.Uid] = om
}

// Get model.OrderModel from cache
func (c Cache) GetOrderById(id string) model.OrderModel {
	return c[id]
}
