package handler

import (
	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/middleware"
	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	cfg      config.App
}

func NewHandler(services *service.Service, cfg config.App) *Handler {
	return &Handler{services: services, cfg: cfg}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(service.CustomLogger(service.LoggerWrite()))
	router.RedirectTrailingSlash = true

	v1 := router.Group("/api/" + h.cfg.Version)
	{
		open := v1.Group("")
		{
			users := open.Group("/users")
			{
				users.POST("/", h.userCreate)
			}

			posts := open.Group("/posts")
			{
				posts.POST("/", h.postCreate)
				posts.GET("/", h.postGetAll)
				posts.GET("/:id", h.postGet)
			}
		}

		auth := v1.Group("", middleware.JWTMiddleware(h.cfg))
		{
			users := auth.Group("/users")
			{
				users.GET("/", middleware.Admin(h.services), h.userGetAll)
				users.GET("/:id", h.userGetByID)
				users.PUT("/:id", h.userUpdate)
				users.DELETE("/:id", h.userDelete)
			}

			posts := auth.Group("/posts")
			{
				posts.PUT("/:id", h.postUpdate)
				posts.DELETE("/:id", h.postDelete)
			}
		}

	}

	return router
}
