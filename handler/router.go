package handler

import (
	"github.com/asliddinberdiev/medium_clone/config"
	"github.com/asliddinberdiev/medium_clone/docs"
	"github.com/asliddinberdiev/medium_clone/middleware"
	"github.com/asliddinberdiev/medium_clone/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.RedirectTrailingSlash = true

	router.Use(service.CustomLogger(service.LoggerWrite()))
	router.Use(middleware.CORS())

	v1 := router.Group("/api/" + h.cfg.Version)
	{
		public := v1.Group("")
		{
			auth := public.Group("/auth")
			{
				auth.POST("/login", h.login)
				auth.POST("/register", h.userCreate)
			}

			posts := public.Group("/posts")
			{
				posts.GET("/", h.postGetAll)
				posts.GET("/:id", h.postGet)
			}

			comments := public.Group("/comments")
			{
				comments.GET("/", h.commentGetAll)
			}
		}

		private := v1.Group("", middleware.JWTMiddleware(h.services))
		{
			auth := private.Group("/auth")
			{
				auth.POST("/logout", h.logout)
			}

			users := private.Group("/users")
			{
				users.GET("/", middleware.Admin(h.services), h.userGetAll)
				users.GET("/:id", h.userGetByID)
				users.PUT("/:id", h.userUpdate)
				users.DELETE("/:id", h.userDelete)
			}

			posts := private.Group("/posts")
			{
				posts.POST("/", h.postCreate)
				posts.GET("/me", h.postGetMe)
				posts.PUT("/:id", h.postUpdate)
				posts.DELETE("/:id", h.postDelete)
			}

			comments := private.Group("/comments")
			{
				comments.POST("/", h.commentCreate)
				comments.PUT("/:id", h.commentUpdate)
				comments.DELETE("/:id", h.commentDelete)
			}

			saved_posts := private.Group("/savedposts")
			{
				saved_posts.POST("/", h.savedPostAdd)
				saved_posts.GET("/", h.savedPostAll)
				saved_posts.DELETE("/:post_id", h.savedPostRemove)
			}
		}
	}

	docs.SwaggerInfo.BasePath = v1.BasePath()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
