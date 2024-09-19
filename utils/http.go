package utils

import (
	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, status int, message string) {
	ctx.AbortWithStatusJSON(status, models.ResponseStatus{Success: false, Message: message})
}

func Data(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.AbortWithStatusJSON(status, models.Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func List(ctx *gin.Context, status int, message string, limit, page uint, data interface{}) {
	ctx.AbortWithStatusJSON(status, models.ResponseList{
		Success: true,
		Message: message,
		Limit:   limit,
		Page:    page,
		Data:    data,
	})
}

func Status(ctx *gin.Context, status int, message string) {
	ctx.AbortWithStatusJSON(status, models.ResponseStatus{Success: true, Message: message})
}
