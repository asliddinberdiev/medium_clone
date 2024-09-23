package handler

import (
	"database/sql"
	"net/http"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

// @Summary      SavedPost
// @Description  save post
// @Tags         saved_post
// @Accept       json
// @Produce      json
// @Param        req  body       models.SavedPostAction true "Save Post Request"
// @Success      200   {object}  models.ResponseStatus
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      403   {object}  models.ResponseStatus
// @Failure      409   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /savedposts [post]
// @Security ApiKeyAuth
func (h *Handler) savedPostAdd(ctx *gin.Context) {
	var input models.SavedPostAction
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	user_id := ctx.GetString("user_id")
	if user_id != input.UserID {
		utils.Error(ctx, http.StatusBadRequest, "user_id wrong")
		return
	}

	_, err := h.services.SavedPost.GetByID(input.UserID, input.PostID)
	if err == nil {
		utils.Status(ctx, http.StatusConflict, "already added")
		return
	} else if err != sql.ErrNoRows {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if err := h.services.SavedPost.Add(input); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Status(ctx, http.StatusOK, "added successfully")
}

// @Summary      SavedPost
// @Description  save post remove
// @Tags         saved_post
// @Accept       json
// @Produce      json
// @Param 		 post_id path string true "post_id"
// @Success      200   {object}  models.ResponseStatus
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      403   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /savedposts/:post_id [delete]
// @Security ApiKeyAuth
func (h *Handler) savedPostRemove(ctx *gin.Context) {
	user_id := ctx.GetString("user_id")
	post_id := ctx.Param("post_id")
	if post_id == "" {
		utils.Error(ctx, http.StatusBadRequest, "required id")
		return
	}

	dbPost, err := h.services.SavedPost.GetByID(user_id, post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Status(ctx, http.StatusNotFound, "not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if user_id != dbPost.UserID {
		utils.Error(ctx, http.StatusForbidden, "access denied")
		return
	}

	if err := h.services.SavedPost.Remove(post_id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Status(ctx, http.StatusOK, "removed successfully")
}

// @Summary      GetAll
// @Description  get all saved post
// @Tags         saved_post
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseList
// @Failure      500  {object}  models.ResponseStatus
// @Router       /posts [get]
func (h *Handler) savedPostAll(ctx *gin.Context) {
	user_id := ctx.GetString("user_id")

	list, err := h.services.SavedPost.GetAll(user_id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.List(ctx, http.StatusOK, "get successfully", uint(len(list)), 1, list)
}
