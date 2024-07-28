package handler

import (
	"SSO-GC/api/request"
	"SSO-GC/config"
	"SSO-GC/internal/app/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

// Handler struct to hold tokens and other related data
type Handler struct {
	cfg *config.AppConfig
	// Optional to hold the SSO tokens
	tokens map[string]interface{}
}

// NewHandler creates a new Handler instance
func NewHandler(cfg *config.AppConfig) *Handler {
	return &Handler{
		cfg:    cfg,
		tokens: make(map[string]interface{}),
	}
}

func (h *Handler) OpenIDConfigHandler(c echo.Context) error {
	config, err := auth.GetOpenIDConfiguration()
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get OpenID Configuration"})
	}
	return c.JSON(http.StatusOK, config)
}

func (h *Handler) TokenHandler(c echo.Context) error {
	var params request.CreateTokensParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request parameters"})
	}

	if code := c.QueryParam("code"); code != "" {
		params.Code = code
	}

	tokens, err := auth.GetTokens(&params)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	// Store the tokens in the Handler struct
	h.tokens = tokens

	return c.JSON(http.StatusOK, tokens)
}

func (h *Handler) UserInfoHandler(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Authorization header is missing"})
	}

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Authorization header format"})
	}
	accessToken := authHeader[len(prefix):]

	userInfo, err := auth.GetUserInfo(accessToken)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user info"})
	}

	return c.JSON(http.StatusOK, userInfo)
}

func (h *Handler) LogoutHandler(c echo.Context) error {
	postLogoutRedirectURI := c.QueryParam("post_logout_redirect_uri")
	idTokenHint := c.QueryParam("id_token_hint")

	// Validate id_token_hint if no cookie is set
	if idTokenHint == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id_token_hint is required"})
	}

	// Remove cookies
	for _, cookie := range c.Cookies() {
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
		c.SetCookie(cookie)
	}

	// If post_logout_redirect_uri is present, respond with 308 and set Location header
	if postLogoutRedirectURI != "" {
		return c.Redirect(http.StatusPermanentRedirect, postLogoutRedirectURI)
	}

	// Otherwise, respond with 200 OK and an empty response
	return c.NoContent(http.StatusOK)
}
