package cache

import (
	"L0/internal/database"
	"L0/internal/models"
	"fmt"
	"sync"
)

type Cache struct {
	db   *database.Postgres
	mu   sync.RWMutex
	data map[string]*models.OrderJSON
}

func NewCache(db *database.Postgres) *Cache {
	return &Cache{
		data: make(map[string]*models.OrderJSON),
		db:   db,
		mu:   sync.RWMutex{},
	}
}
func (c *Cache) GetOrder(uid string) *models.OrderJSON {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[uid]
}
func (c *Cache) GetAllOrders() map[string]*models.OrderJSON {
	return c.data
}
func (c *Cache) AddCache(order models.OrderJSON) {
	err := c.db.SaveOrder(order)
	if err != nil {
		fmt.Printf("Cannot insert order: %v\n", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[order.OrderUid] = &order
	fmt.Printf("Cache written: %s\n", order.OrderUid)
}
func (c *Cache) Preload() error {
	array, err := c.db.GetAllOrders()
	if err != nil {
		return err
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, order := range array {
		c.data[order.OrderUid] = &order
	}
	return nil
}
