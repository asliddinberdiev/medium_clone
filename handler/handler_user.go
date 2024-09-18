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

	accessToken, err := h.services.Token.AccessTokenGenerate(user.ID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	refreshToken, err := h.services.Token.RefreshTokenGenerate(user.ID)
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

func (h *Handler) userGetAll(ctx *gin.Context) {
	list, err := h.services.User.GetAll()
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "get all user successfully", list)
}

func (h *Handler) userGetByID(ctx *gin.Context) {
	id, ok := ctx.Params.Get("id")
	if !ok {
		utils.Error(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.User.GetByID(id)
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

	id, ok := ctx.Params.Get("id")
	if !ok {
		utils.Error(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.User.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(ctx, http.StatusNotFound, "user not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if input.Role != "" && user.Role != "admin" {
		utils.Error(ctx, http.StatusForbidden, "access denied")
		return
	}

	updateUser, err := h.services.User.Update(id, input)
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
	id, ok := ctx.Params.Get("id")
	if !ok {
		utils.Error(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	err := h.services.User.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(ctx, http.StatusNotFound, "user not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "delete user successfully", map[string]interface{}{
		"id": id,
	})
}
