package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"user_service/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	svc       *service.UserService
	jwtSecret string
}

func NewUserHandler(svc *service.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{svc: svc, jwtSecret: jwtSecret}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}
	user, err := h.svc.Register(req.Username, req.Password, req.FullName)

	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	token, user, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": token,
		"token_type":   "Bearer",
		"user_id":      user.ID,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	idStr := c.Request().Header.Get("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or missing X-User-Id header"})
	}

	user, err := h.svc.GetProfile(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUploadURL(c echo.Context) error {
	idStr := c.Request().Header.Get("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or missing X-User-Id header"})
	}

	fileName := c.QueryParam("file_name")
	if fileName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid or missing file_name parameter"})
	}

	url, err := h.svc.GenerateUploadURL(uint(id), fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"upload_url":         url,
		"expires_in_seconds": 300,
		"method":             "PUT",
	})
}

func (h *UserHandler) VerifyToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing token"})
	}

	tokenStr := authHeader[7:]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.ErrUnauthorized
		}
		return []byte(h.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%.0f", claims["user_id"].(float64))
	c.Response().Header().Set("X-User-Id", userID)
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) Dashboard(c echo.Context) error {
	idStr := c.Request().Header.Get("X-User-Id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing X-User-Id"})
	}
	result, err := h.svc.Dashboard(uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
