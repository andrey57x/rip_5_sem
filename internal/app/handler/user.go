package handler

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) SignIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_in",
	})
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByID(id)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

func (h *Handler) ChangeProfile(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.ChangeProfile(id, userJSON)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

func (h *Handler) SignOut(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_out",
	})
}
