package handler

import (
	apitypes "Backend/internal/app/api_types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateUser godoc
// @Summary Create user
// @Description Зарегистрироваться
// @Tags users
// @Accept json
// @Produce json
// @Param user body apitypes.UserJSON true "Create user"
// @Success 201 {object} apitypes.UserJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/sign-up [post]
func (h *Handler) CreateUser(ctx *gin.Context) {
	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.CreateUser(userJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/users/%v", user.ID))
	ctx.JSON(http.StatusCreated, apitypes.UserToJSON(user))
}

// SignIn godoc
// @Summary Sign in
// @Description Войти
// @Tags users
// @Accept json
// @Produce json
// @Param user body apitypes.UserJSON true "Sign in"
// @Success 200 {object} apitypes.UserJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/sign-in [post]
func (h *Handler) SignIn(ctx *gin.Context) {
	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.SignIn(userJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

// GetProfile godoc
// @Summary Get profile
// @Description Получить личный кабинет
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} apitypes.UserJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [get]
func (h *Handler) GetProfile(ctx *gin.Context) {
	user, err := h.Repository.GetUserByID(h.Repository.GetUserID())
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

// ChangeProfile godoc
// @Summary Change profile
// @Description Изменить личный кабинет
// @Tags users
// @Accept json
// @Produce json
// @Param user body apitypes.UserJSON true "Change profile"
// @Success 200 {object} apitypes.UserJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [put]
func (h *Handler) ChangeProfile(ctx *gin.Context) {
	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	user, err := h.Repository.ChangeProfile(h.Repository.GetUserID(), userJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

// SignOut godoc
// @Summary Sign out
// @Description Выйти
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/sign-out [post]
func (h *Handler) SignOut(ctx *gin.Context) {
	h.Repository.SignOut()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_out",
	})
}
