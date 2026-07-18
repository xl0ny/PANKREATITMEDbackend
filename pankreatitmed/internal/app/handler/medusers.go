package handler

import (
	"errors"
	"net/http"
	"pankreatitmed/internal/app/authctx"
	"pankreatitmed/internal/app/dto/request"
	"pankreatitmed/internal/app/dto/response"
	"pankreatitmed/internal/app/services"

	"github.com/gin-gonic/gin"
)

// MedUserRegistation godoc
// @Summary      Регистрация
// @Description  Регистрирует нового пользователя и сразу возвращает JWT (Bearer)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body request.MedUserRegistration true "Логин/пароль и прочие поля"
// @Success      201 {object} response.AuthorizateUser
// @Failure      400 {object} map[string]any "bad request / weak password"
// @Failure      409 {object} map[string]any "login already taken"
// @Failure      500 {object} map[string]any "internal error"
// @Router       /users/auth/register [post]
func (h *Handler) MedUserRegistation(c *gin.Context) {
	var user request.MedUserRegistration
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, token, err := h.svcs.MedUsers.Register(user)
	switch {
	case errors.Is(err, services.ErrWeakPassword):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	case errors.Is(err, services.ErrLoginTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	exp := h.svcs.MedUsers.GetConfig().TTL.Hours()
	c.JSON(http.StatusCreated, response.AuthorizateUser{AccessToken: token, TokenType: "Bearer", ExpiresIn: int(exp)})
}

// MedUserGetFields godoc
// @Summary      Личный кабинет: получить мои поля
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} response.SendMedUserField
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      404 {object} map[string]any "not found"
// @Router       /users/me [get]
func (h *Handler) MedUserGetFields(c *gin.Context) {
	user, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	res, err := h.svcs.MedUsers.GetMyField(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// MedUserUpdateFields godoc
// @Summary      Личный кабинет: обновить мои поля
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        input body request.UpdateMedUser true "Изменяемые поля"
// @Success      200 {string} string "ok"
// @Failure      400 {object} map[string]any "bad request / validation error"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      404 {object} map[string]any "not found"
// @Router       /users/me [put]
func (h *Handler) MedUserUpdateFields(c *gin.Context) {
	usr, check := authctx.Get(c)
	if !check {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "problem with your token"})
		return
	}
	var user request.UpdateMedUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svcs.MedUsers.UpdateField(usr.ID, &user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// MedUserLogIn godoc
// @Summary      Вход (логин)
// @Description  Проверяет логин/пароль и возвращает JWT (Bearer)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body request.AuthenticateMedUser true "Логин и пароль"
// @Success      200 {object} response.AuthorizateUser
// @Failure      400 {object} map[string]any "bad request"
// @Failure      401 {object} map[string]any "invalid credentials"
// @Router       /users/auth/login [post]
func (h *Handler) MedUserLogIn(c *gin.Context) {
	var acces request.AuthenticateMedUser
	if err := c.ShouldBindJSON(&acces); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	token, err := h.svcs.MedUsers.Login(acces)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	exp := h.svcs.MedUsers.GetConfig().TTL.Hours()
	c.JSON(http.StatusOK, response.AuthorizateUser{AccessToken: token, TokenType: "Bearer", ExpiresIn: int(exp)})
}

// MedUserLogOut godoc
// @Summary      Выход (logout) — поместить токен в blacklist
// @Description  Добавляет переданный токен в blacklist до конца срока его действия
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Param        token path string true "JWT или jti (в зависимости от реализации Logout)"
// @Success      200 {object} map[string]any "status/message"
// @Failure      401 {object} map[string]any "unauthenticated"
// @Failure      403 {object} map[string]any "forbidden"
// @Failure      404 {object} map[string]any "token not found / already revoked"
// @Router       /users/auth/logout/{token} [post]
func (h *Handler) MedUserLogOut(c *gin.Context) {
	token := c.Param("token")
	if err := h.svcs.MedUsers.Logout(token); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "logout success",
	})
}
