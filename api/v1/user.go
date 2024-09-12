package v1

import (
	"database/sql"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/asliddinberdiev/medium_clone/api/models"
	"github.com/asliddinberdiev/medium_clone/storage/repo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handlerV1) CreateUser(ctx *gin.Context) {
	var req models.CreateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err,
		})
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("user create uuid error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	user, err := h.strg.User().Create(ctx, &repo.User{
		ID:        id.String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {

		if strings.Contains(err.Error(), "unique") {
			log.Printf("this email already used error: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "This email already used",
			})
			return
		}

		log.Printf("user create storage error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "ok",
		"message": "User created successfully",
		"user_id": user.ID,
	})
}

func (h *handlerV1) GetUser(ctx *gin.Context) {
	user_id := ctx.Param("id")
	if user_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user id requered",
		})
		return
	}

	user, err := h.strg.User().Get(ctx, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		log.Println("get user error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Get user successfully",
		"user": models.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func (h *handlerV1) UpdateUser(ctx *gin.Context) {
	user_id := ctx.Param("id")
	if user_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user id requered",
		})
		return
	}

	var req models.UpdateUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "a field requered",
		})
		return
	}

	if reflect.DeepEqual(req, models.UpdateUser{}) {
		log.Println("user update data empty")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "a field requered",
		})
		return
	}

	user, err := h.strg.User().Get(ctx, user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		log.Println("get user error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	if req.FirstName == "" {
		req.FirstName = user.FirstName
	}

	if req.LastName == "" {
		req.LastName = user.LastName
	}

	if req.Password == "" {
		req.Password = user.Password
	}

	if err := h.strg.User().Update(ctx, &repo.UpdateUser{
		ID:        user_id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
	}); err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		log.Println("get user error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User updated successfully",
	})
}

func (h *handlerV1) DeleteUser(ctx *gin.Context) {
	user_id := ctx.Param("id")
	if user_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "user id requered",
		})
		return
	}

	if err := h.strg.User().Delete(ctx, user_id); err != nil {
		if err == sql.ErrNoRows {
			log.Println("user not found: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}

		log.Println("delete user error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal error we got :(",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User deleted successfully",
	})
}
