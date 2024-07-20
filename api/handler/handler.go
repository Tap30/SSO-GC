package handler

import (
	"SSO-GC/api/request"
	"SSO-GC/internal/app/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func OpenIDConfigHandler(c echo.Context) error {
	config, err := auth.GetOpenIDConfiguration()
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get OpenID Configuration"})
	}
	return c.JSON(http.StatusOK, config)
}

func TokenHandler(c echo.Context) error {
	var params request.CreateTokensParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request parameters"})
	}

	tokens, err := auth.GetTokens(&params)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	return c.JSON(http.StatusOK, tokens)
}
