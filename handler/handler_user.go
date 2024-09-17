package handler

import (
	"net/http"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) userCreate(ctx *gin.Context) {
	var input models.UserCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "ShouldBindJSON")
		return
	}

	user, err := h.services.User.Create(ctx, input)
	if err != nil {
		if err.Error() == "unique" {
			utils.Error(ctx, http.StatusBadRequest, "this email already used")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	accessToken, err := h.services.Token.AccessTokenGenerate(user.ID, user.Role)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	refreshToken, err := h.services.Token.RefreshTokenGenerate(user.ID, user.Role)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	token := map[string]string{
		"access":  accessToken,
		"refresh": refreshToken,
	}

	utils.Data(ctx, http.StatusCreated, "created successfully", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
func (h *Handler) userGet(c *gin.Context)    {}
func (h *Handler) userUpdate(c *gin.Context) {}
func (h *Handler) userDelete(c *gin.Context) {}
