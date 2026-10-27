package handler

import (
	"net/http"
	"order_service/internal/service"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	idStr := c.Request().Header.Get("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing X-User-Id"})
	}

	var req struct {
		ItemName string  `json:"item_name"`
		Amount   float64 `json:"amount"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	order, err := h.svc.CreateOrder(uint(id), req.ItemName, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	idStr := c.Request().Header.Get("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing X-User-Id"})
	}

	orders, err := h.svc.GetOrdersByUser(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, orders)
}
