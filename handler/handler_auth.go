package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/asliddinberdiev/medium_clone/models"
	"github.com/asliddinberdiev/medium_clone/utils"
	"github.com/gin-gonic/gin"
)

// @Summary      Login
// @Description  login a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.Login true "Login Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /auth/login [post]
func (h *Handler) login(ctx *gin.Context) {
	var input models.Login
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if err := input.IsValid(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "fields validation error")
		return
	}

	user, err := h.services.User.GetByEmail(input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.Error(ctx, http.StatusBadRequest, "email or password wrong")
			return
		}
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if ok := utils.CheckPassword(user.Password, input.Password); !ok {
		utils.Error(ctx, http.StatusBadRequest, "email or password wrong")
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

	utils.Data(ctx, http.StatusOK, "welcome", map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// @Summary      Logout
// @Description  logout a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req  body      models.Token true "Logout Request"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.ResponseStatus
// @Failure      400  {object}  models.ResponseStatus
// @Failure      401  {object}  models.ResponseStatus
// @Failure      500  {object}  models.ResponseStatus
// @Router       /auth/logout [post]
// @Security ApiKeyAuth
func (h *Handler) logout(ctx *gin.Context) {
	var input models.Token
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "required fields")
		return
	}

	if input.Token == "" {
		utils.Error(ctx, http.StatusBadRequest, "invalid token 1")
		return
	}

	claims, err := h.services.Token.Parse(input.Token)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		log.Println("handler_auth: logout - claims  jti type error")
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	if err := h.services.Auth.AddBlack(tokenID, input.Token, time.Hour*24); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "we got internal server :(")
		return
	}

	utils.Status(ctx, http.StatusOK, "see you again :)")
}
