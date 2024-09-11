package v1

import (
	"log"
	"net/http"
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

func (h *handlerV1) GetUser(ctx *gin.Context) {

}

func (h *handlerV1) UpdateUser(ctx *gin.Context) {

}

func (h *handlerV1) DeleteUser(ctx *gin.Context) {

}
