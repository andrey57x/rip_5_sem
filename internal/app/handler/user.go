package handler

import (
	apitypes "Backend/internal/app/api_types"
	"fmt"
	"net/http"

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

func (h *Handler) GetProfile(ctx *gin.Context) {
	user, err := h.Repository.GetUserByID(h.Repository.GetUserID())
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

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

func (h *Handler) SignOut(ctx *gin.Context) {
	h.Repository.SignOut()
	ctx.JSON(http.StatusOK, gin.H{
		"status": "signed_out",
	})
}
