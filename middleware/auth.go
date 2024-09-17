package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenBearer := c.GetHeader("Authorization")

		log.Printf("token: %v", tokenBearer)

		c.Next()
	}
}
