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
		ctx.Set("role", claims["role"])

		ctx.Next()
	}
}
