package handler

import (
	"database/sql"
	"net/http"
	"reflect"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handler) userCreate(ctx *gin.Context) {
	var input models.UserCreate
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	user, err := h.services.User.Create(input)
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

func (h *Handler) userGet(ctx *gin.Context) {
	user_id := ctx.GetString("user_id")
	id, ok := ctx.Params.Get("id")
	if !ok {
		utils.Error(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	if user_id != id {
		utils.Error(ctx, http.StatusNotFound, "invalid id param")
		return
	}

	user, err := h.services.User.GetByID(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(ctx, http.StatusNotFound, "user not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "get user successfully", user)
}

func (h *Handler) userUpdate(ctx *gin.Context) {
	var input models.UpdateUser
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}
	if reflect.DeepEqual(input, models.UpdateUser{}) {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	user_id := ctx.GetString("user_id")
	id, ok := ctx.Params.Get("id")
	if !ok {
		utils.Error(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	if user_id != id {
		utils.Error(ctx, http.StatusNotFound, "invalid id param")
		return
	}

	updateUser, err := h.services.User.Update(user_id, input)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(ctx, http.StatusNotFound, "user not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "get user successfully", updateUser)
}

func (h *Handler) userDelete(ctx *gin.Context) {
	user_id := ctx.GetString("user_id")
	id, ok := ctx.Params.Get("id")
	if !ok {
		utils.Error(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	if user_id != id {
		utils.Error(ctx, http.StatusNotFound, "invalid id param")
		return
	}

	err := h.services.User.Delete(user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(ctx, http.StatusNotFound, "user not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "delete user successfully", map[string]interface{}{
		"id": user_id,
	})
}
