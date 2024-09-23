package handler

import (
	"database/sql"
	"net/http"
	"reflect"

	models "github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

// @Summary      Register
// @Description  register a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.UserCreate true "Register Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /auth/register [post]
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
		if utils.HasStringKey(err.Error(), "duplicate key") {
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

// @Summary      GetAll
// @Description  get all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200   {object}   models.ResponseList
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      403   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /users [get]
// @Security ApiKeyAuth
func (h *Handler) userGetAll(ctx *gin.Context) {
	list, err := h.services.User.GetAll()
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.List(ctx, http.StatusOK, "get all user successfully", 10, 1, list)
}

// @Summary      GetByID
// @Description  get by id user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200   {object}   models.Response
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      404   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /users/{id} [get]
// @Security ApiKeyAuth
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

// @Summary      Update
// @Description  update user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Param 		 req  body  models.UpdateUser true "UpdateUser Request"
// @Success      200   {object}   models.Response
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      403   {object}  models.ResponseStatus
// @Failure      404   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /users/{id} [put]
// @Security ApiKeyAuth
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

	utils.Data(ctx, http.StatusOK, "update user successfully", updateUser)
}

// @Summary      Delete
// @Description  delete user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200   {object}   models.Response
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      404   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /users/{id} [delete]
// @Security ApiKeyAuth
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
