package httpapi

import (
	"user_service/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(h *handler.UserHandler) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("/users")
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.GET("/profile", h.GetProfile)
	g.GET("/avatar-url", h.GetUploadURL)

	return e
}
