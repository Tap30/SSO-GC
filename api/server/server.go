package server

import (
	"SSO-GC/api/handler"
	"SSO-GC/config"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Allows all origins
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	h := handler.NewHandler()

	e.GET("/.well-known/openid-configuration", h.OpenIDConfigHandler)
	e.POST("/token", h.TokenHandler)
	e.POST("/userinfo", h.UserInfoHandler)
	e.GET("/logout", h.LogoutHandler)

	cfg := config.LoadConfig()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.ServerPort)))
}
