package router

import (
	v1 "github.com/asliddinberdiev/medium_clone/api/v1"
	"github.com/asliddinberdiev/medium_clone/storage"
	"github.com/gin-gonic/gin"
)

type Options struct {
	Strg storage.StorageI
}

func NewRouter(opts *Options) *gin.Engine {
	router := gin.New()

	// ROUTER
	handler := v1.New(&v1.HandlerV1{Strg: opts.Strg})

	// user routes
	router.POST("/v1/users", handler.CreateUser)
	router.GET("/v1/users/:id", handler.GetUser)
	router.PUT("/v1/users/:id", handler.UpdateUser)
	router.DELETE("/v1/users/:id", handler.DeleteUser)

	return router
}
