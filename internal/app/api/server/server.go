package server

import (
	"SSO-GC/internal/app/api/handler"
	"SSO-GC/internal/app/api/server/custommiddleware"
	"SSO-GC/internal/app/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// StartServer initializes and starts the Echo server with configured routes and middleware.
func StartServer() {
	cfg := config.LoadConfig()
	h := handler.NewHandler(cfg)

	e := echo.New()
	e.Use(middleware.CORS())

	// Use custom KeyAuth middleware
	keyAuthMiddleware := custommiddleware.KeyAuthMiddleware()

	e.GET("/.well-known/openid-configuration", h.OpenIDConfigHandler)
	e.POST("/token", h.TokenHandler)
	e.POST("/userinfo", h.UserInfoHandler, keyAuthMiddleware)
	e.GET("/logout", h.LogoutHandler, keyAuthMiddleware)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.ServerPort)))
	// sd
}
