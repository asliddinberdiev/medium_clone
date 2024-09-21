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
// @Description  create a post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        req  body      models.CreatePost true "Create Post Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /posts [post]
// @Security ApiKeyAuth
func (h *Handler) postCreate(ctx *gin.Context) {
	var input models.CreatePost
	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	user_id := ctx.GetString("user_id")

	post, err := h.services.Post.Create(user_id, input)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "created successfully", post)
}

// @Summary      GetAll
// @Description  get all post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseList
// @Failure      400  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /posts [get]
func (h *Handler) postGetAll(ctx *gin.Context) {
	list, err := h.services.Post.GetAll()
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.List(ctx, http.StatusOK, "get all successfully", uint(len(list)), 1, list)
}

// @Summary      GetMe
// @Description  get all my posts
// @Tags         posts
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.ResponseList
// @Failure      400  {object}  models.ResponseStatus
// @Failure      401  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /posts [get]
// @Security ApiKeyAuth
func (h *Handler) postGetMe(ctx *gin.Context) {
	user_id := ctx.GetString("user_id")
	list, err := h.services.Post.GetPersonal(user_id)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.List(ctx, http.StatusOK, "get all successfully", uint(len(list)), 1, list)
}

// @Summary      GetByID
// @Description  get by id post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /posts/{id} [get]
func (h *Handler) postGet(ctx *gin.Context) {
	post_id := ctx.Param("id")

	post, err := h.services.Post.GetByID(post_id)
	if err != err {
		if err == sql.ErrNoRows {
			utils.Status(ctx, http.StatusNotFound, "not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "get successfully", post)
}

// @Summary      Update
// @Description  update post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      403  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /posts/{id} [put]
// @Security ApiKeyAuth
func (h *Handler) postUpdate(ctx *gin.Context) {
	var input models.UpdatePost
	post_id := ctx.Param("id")
	if post_id == "" {
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

	dbpost, err := h.services.Post.GetByID(post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Status(ctx, http.StatusNotFound, "not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	user_id := ctx.GetString("user_id")
	if user_id != dbpost.UserID {
		utils.Error(ctx, http.StatusForbidden, "access denied")
		return
	}

	post, err := h.services.Post.Update(post_id, input)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "update successfully", post)
}

// @Summary      Delete
// @Description  delete post
// @Tags         post
// @Accept       json
// @Produce      json
// @Param 		 id path string true "id"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.ResponseStatus
// @Failure      401   {object}  models.ResponseStatus
// @Failure      403   {object}  models.ResponseStatus
// @Failure      404   {object}  models.ResponseStatus
// @Failure      500   {object}  models.ResponseStatus
// @Router       /posts/{id} [delete]
// @Security ApiKeyAuth
func (h *Handler) postDelete(ctx *gin.Context) {
	post_id := ctx.Param("id")
	user_id := ctx.GetString("user_id")

	if post_id == "" {
		utils.Error(ctx, http.StatusBadRequest, "required id")
		return
	}

	post, err := h.services.Post.GetByID(post_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Status(ctx, http.StatusNotFound, "not found")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if user_id != post.UserID {
		utils.Error(ctx, http.StatusForbidden, "access denied")
		return
	}

	if err := h.services.Post.Delete(post_id); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Data(ctx, http.StatusOK, "delete post successfully", map[string]interface{}{
		"id": post_id,
	})
}
