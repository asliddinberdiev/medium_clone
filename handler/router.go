package handler

import (
	"github.com/asliddinberdiev/medium_clone/middleware"
	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	version  string
}

func NewHandler(services *service.Service, version string) *Handler {
	return &Handler{services: services, version: version}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.LoggerWithWriter(service.LoggerWrite()))

	router.RedirectTrailingSlash = true

	// ROUTES
	v1 := router.Group("/api/" + h.version)
	{
		users := v1.Group("/users", middleware.JWTMiddleware())
		{
			users.POST("/", h.userCreate)
			users.GET("/:id", h.userGet)
			users.PUT("/:id", h.userUpdate)
			users.DELETE("/:id", h.userDelete)
		}

		posts := v1.Group("/posts")
		{
			posts.POST("/", h.postCreate)
			posts.GET("/", h.postGetAll)
			posts.GET("/:id", h.postGet)
			posts.PUT("/:id", h.postUpdate)
			posts.DELETE("/:id", h.postDelete)
		}
	}

	return router
}
