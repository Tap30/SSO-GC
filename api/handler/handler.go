package handler

import (
	"net/http"
	"time"

	"SSO-GC/api/request"
	"SSO-GC/config"
	"SSO-GC/internal/app/auth"
	"github.com/labstack/echo/v4"
)

// Handler manages tokens and OpenID configurations.
type Handler struct {
	cfg    *config.AppConfig
	tokens map[string]interface{} // Optional
}

// NewHandler initializes a new Handler instance.
func NewHandler(cfg *config.AppConfig) *Handler {
	return &Handler{
		cfg:    cfg,
		tokens: make(map[string]interface{}),
	}
}

// OpenIDConfigHandler fetches the OpenID configuration.
func (h *Handler) OpenIDConfigHandler(c echo.Context) error {
	openIDConfiguration, err := auth.GetOpenIDConfiguration()
	if err != nil {
		c.Logger().Errorf("Failed to get OpenID Configuration: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get OpenID Configuration"})
	}
	return c.JSON(http.StatusOK, openIDConfiguration)
}

// TokenHandler generates tokens based on the request parameters.
func (h *Handler) TokenHandler(c echo.Context) error {
	var params request.CreateTokensParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request parameters"})
	}

	tokens, err := auth.GetTokens(&params)
	if err != nil {
		c.Logger().Errorf("Failed to generate tokens: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate tokens"})
	}

	// Set cookies for the tokens
	if accessToken, ok := tokens["access_token"].(string); ok {
		c.SetCookie(&http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
	}

	if refreshToken, ok := tokens["refresh_token"].(string); ok {
		c.SetCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
		})
	}

	if idToken, ok := tokens["id_token"].(string); ok {
		c.SetCookie(&http.Cookie{
			Name:  "id_token",
			Value: idToken,
		})
	}

	h.tokens = tokens
	return c.JSON(http.StatusOK, tokens)
}

// UserInfoHandler retrieves user information using the access token.
func (h *Handler) UserInfoHandler(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Authorization header is missing"})
	}

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Authorization header format"})
	}

	userInfo, err := auth.GetUserInfo(authHeader[len(prefix):])
	if err != nil {
		c.Logger().Errorf("Failed to get user info: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user info"})
	}

	return c.JSON(http.StatusOK, userInfo)
}

// LogoutHandler logs out the user by clearing cookies and managing redirect.
func (h *Handler) LogoutHandler(c echo.Context) error {
	postLogoutRedirectURI := c.QueryParam("post_logout_redirect_uri")
	idTokenHint := c.QueryParam("id_token_hint")

	if idTokenHint == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id_token_hint is required"})
	}

	for _, cookie := range c.Cookies() {
		cookie.MaxAge = -1
		cookie.Expires = time.Unix(0, 0)
		c.SetCookie(cookie)
	}

	if postLogoutRedirectURI != "" {
		c.Response().Header().Set("Location", postLogoutRedirectURI)
		return c.JSON(http.StatusPermanentRedirect, nil)
	}

	return c.NoContent(http.StatusOK)
}
