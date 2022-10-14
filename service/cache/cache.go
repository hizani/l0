package cache

import (
	"wbintern/l0/service/database"
	"wbintern/l0/service/model"
)

// Cache for orders
type Cache struct {
	Data map[string]model.OrderModel
}

// Initialize cache
func New() Cache {
	var c Cache
	c.Data = make(map[string]model.OrderModel)
	return c
}

// Restore cache from a database
func (c *Cache) Restore(db *database.Database) error {
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
func (c *Cache) Add(om model.OrderModel) {
	c.Data[om.Uid] = om
}
