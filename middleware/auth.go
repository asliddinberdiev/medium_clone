package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWTMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenBearer := c.GetHeader("Authorization")

		log.Info("token", zap.Any("bearer", tokenBearer))

		c.Next()
	}
}
