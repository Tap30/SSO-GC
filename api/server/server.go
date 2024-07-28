package server

import (
	"SSO-GC/api/handler"
	"SSO-GC/config"
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

	e.GET("/.well-known/openid-configuration", h.OpenIDConfigHandler)
	e.POST("/token", h.TokenHandler)
	e.POST("/userinfo", h.UserInfoHandler)
	e.GET("/logout", h.LogoutHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.ServerPort)))
}
