package middleware

import (
	"net/http"
	"strings"

	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(cfg config.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(ctx, http.StatusUnauthorized, "invalid authorization")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.Error(ctx, http.StatusUnauthorized, "invalid authorization")
			return
		}

		token := tokenParts[1]
		if token == "" {
			utils.Error(ctx, http.StatusUnauthorized, "invalid authorization")
			return
		}

		claims, err := service.NewTokenService(cfg).Parse(token)
		if err != nil {
			utils.Error(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set("user_id", claims["id"])

		ctx.Next()
	}
}

func Admin(services *service.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString("user_id")

		user, err := services.User.GetByID(id)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
			return
		}

		if user.Role != "admin" {
			utils.Error(ctx, http.StatusForbidden, "access denied")
			return
		}

		ctx.Next()
	}
}
