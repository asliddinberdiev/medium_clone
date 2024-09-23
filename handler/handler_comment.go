package handler

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

// @Summary      Create
// @Description  create a comment
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        req  body      models.CreateComment true "Create Comment Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      401  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /comments [post]
// @Security ApiKeyAuth
func (h *Handler) commentCreate(ctx *gin.Context) {
	var input models.CreateComment
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	user_id := ctx.GetString("user_id")

	_, err := h.services.Post.GetByID(input.PostID)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "post not found")
		return
	}

	comment, err := h.services.Comment.Create(user_id, input)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "created successfully", comment)
}

// @Summary      GetAll
// @Description  get all comments
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param        post_id query string false "Post ID"
// @Success      200  {object}  models.ResponseList
// @Failure      400  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /comments [get]
func (h *Handler) commentGetAll(ctx *gin.Context) {
	post_id := ctx.Query("post_id")
	if post_id == "" {
		utils.Error(ctx, http.StatusBadRequest, "required query post_id")
		return
	}

	list, err := h.services.Comment.GetAll(post_id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.List(ctx, http.StatusOK, "get successfully", uint(len(list)), 1, list)
}

// @Summary      Update
// @Description  update comment
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Param        req  body      models.UpdateComment true "Update Comment Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      401  {object}  models.ResponseStatus
// @Failure      403  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /comments/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) commentUpdate(ctx *gin.Context) {
	var input models.UpdateComment
	comment_id := ctx.Param("id")
	if comment_id == "" {
		utils.Error(ctx, http.StatusBadRequest, "required id")
		return
	}

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required field")
		return
	}

	if reflect.DeepEqual(input, models.UpdateUser{}) {
		utils.Error(ctx, http.StatusBadRequest, "required a field")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	dbComment, err := h.services.Comment.GetByID(comment_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Status(ctx, http.StatusNotFound, "not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	user_id := ctx.GetString("user_id")
	if user_id != dbComment.UserID {
		utils.Error(ctx, http.StatusForbidden, "access denied")
		return
	}

	comment, err := h.services.Comment.Update(comment_id, input.Body)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "update successfully", comment)
}

// @Summary      Delete
// @Description  delete comment
// @Tags         comment
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      403   {object}  models.ResponseStatus
// @Failure      404   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /comments/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) commentDelete(ctx *gin.Context) {
	comment_id := ctx.Param("id")
	user_id := ctx.GetString("user_id")

	if comment_id == "" {
		utils.Error(ctx, http.StatusBadRequest, "required id")
		return
	}

	comment, err := h.services.Comment.GetByID(comment_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Status(ctx, http.StatusNotFound, "not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if user_id != comment.UserID {
		utils.Error(ctx, http.StatusForbidden, "access denied")
		return
	}

	if err := h.services.Comment.Delete(comment_id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "delete comment successfully", map[string]interface{}{
		"id": comment_id,
	})
}
