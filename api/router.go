package router

import (
	v1 "github.com/asliddinberdiev/medium_clone/api/v1"
	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/storage"
	"github.com/gin-gonic/gin"
)

type Options struct {
	Strg storage.StorageI
}

func NewRouter(opts *Options) *gin.Engine {
	cfg := config.Load(".")
	router := gin.New()
	router.RedirectTrailingSlash = true

	// ROUTER
	handler := v1.New(&v1.HandlerV1{Strg: opts.Strg})

	v1 := router.Group("/api/" + cfg.App.Version)
	{
		users := v1.Group("/users")
		{
			users.POST("/", handler.CreateUser)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

		posts := v1.Group("/posts")
		{
			posts.POST("/", handler.CreatePost)
			posts.GET("/:id", handler.GetPost)
			posts.PUT("/:id", handler.UpdatePost)
			posts.DELETE("/:id", handler.DeletePost)
		}
	}

	return router
}
