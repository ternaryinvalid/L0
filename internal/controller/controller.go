package controller

import (
	"L0/internal/cache"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CacheHandler struct {
	c *cache.Cache
}

func NewOrderController(c *cache.Cache) *CacheHandler {
	return &CacheHandler{
		c: c,
	}
}

func (c *CacheHandler) GetOrder(ctx echo.Context) error {
	order := c.c.GetOrder(ctx.Param("order"))

	ord, err := json.MarshalIndent(order, "", "\t")

	if err != nil {
		return err
	}

	return ctx.JSONBlob(http.StatusOK, ord)
}

func (c *CacheHandler) GetAllOrder(ctx echo.Context) error {
	order := c.c.GetAllOrders()

	ord, err := json.MarshalIndent(order, "", "\t")

	if err != nil {
		return err
	}

	return ctx.JSONBlob(http.StatusOK, ord)

}
