package httpapi

import (
	"order_service/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(h *handler.OrderHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("/orders")
	g.POST("", h.CreateOrder)
	g.GET("", h.GetOrders)

	return e
}
