package v1

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/asliddinberdiev/medium_clone/api/models"
	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handlerV1) CreatePost(ctx *gin.Context) {
	var req models.CreatePost
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err,
		})
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("post create uuid err: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	post, err := h.strg.Post().Create(ctx, &repo.Post{
		ID:        id.String(),
		UserID:    req.UserID,
		Title:     req.Title,
		Body:      req.Body,
		Published: req.Published,
	})
	if err != nil {
		log.Println("post create storage err: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Post created successfully",
		"post": models.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Body:      post.Body,
			Published: post.Published,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func (h *handlerV1) GetPost(ctx *gin.Context) {
	post_id := ctx.Param("id")
	if post_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "post id requered",
		})
		return
	}

	post, err := h.strg.Post().Get(ctx, post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("post not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "post not found",
			})
			return
		}
		log.Println("get post error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Get post successfully",
		"post": models.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Body:      post.Body,
			Published: post.Published,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
			UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func (h *handlerV1) UpdatePost(ctx *gin.Context) {
	post_id := ctx.Param("id")
	if post_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "post id requered",
		})
		return
	}

	var req models.UpdatePost
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "a field requered",
		})
		return
	}

	post, err := h.strg.Post().Get(ctx, post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("post not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "post not found",
			})
			return
		}

		log.Println("get post error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	if req.Title == "" {
		req.Title = post.Title
	}

	if req.Body == "" {
		req.Body = post.Body
	}

	if !req.Published {
		req.Published = post.Published
	}

	if err := h.strg.Post().Update(ctx, &repo.UpdatePost{
		ID:        post_id,
		Title:     req.Title,
		Body:      req.Body,
		Published: req.Published,
	}); err != nil {
		if err == sql.ErrNoRows {
			log.Println("post not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "post not found",
			})
			return
		}

		log.Println("update post error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Post updated successfully",
	})
}

func (h *handlerV1) DeletePost(ctx *gin.Context) {
	post_id := ctx.Param("id")
	if post_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "post id requered",
		})
		return
	}

	if err := h.strg.Post().Delete(ctx, post_id); err != nil {
		if err == sql.ErrNoRows {
			log.Println("post not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "post not found",
			})
			return
		}

		log.Println("delete post error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Post deleted successfully",
	})
}
