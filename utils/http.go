package utils

import (
	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

func Data(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func List(ctx *gin.Context, status int, message string, limit, page uint, data []interface{}) {
	ctx.JSON(status, gin.H{
		"success": true,
		"message": message,
		"limit":   limit,
		"page":    page,
		"data":    data,
	})
}

func Status(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, gin.H{
		"success": true,
		"message": message,
	})
}
