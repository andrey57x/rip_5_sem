package handler

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/repository"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) SignUp(ctx *gin.Context) {
	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.CreateUser(userJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/users/%v", user.Login))
	ctx.JSON(http.StatusCreated, apitypes.UserToJSON(user))
}

func (h *Handler) SignIn(ctx *gin.Context) {
	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	token, err := h.Repository.SignIn(userJSON)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Handler) GetProfile(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	login := ctx.Param("login")

	user, err := h.Repository.GetUserByLogin(login)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	if user.UUID != userID {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("users do not match"))
		return
	}

	ctx.JSON(http.StatusOK, apitypes.UserToJSON(user))
}

func (h *Handler) ChangeProfile(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	login := ctx.Param("login")

	var userJSON apitypes.UserJSON
	if err := ctx.BindJSON(&userJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByLogin(login)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if user.UUID != userID {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}

	user, err = h.Repository.ChangeProfile(login, userJSON)
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
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.SignOut(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error deleting token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "signed_out"})
}

func getUserID(ctx *gin.Context) (uuid.UUID, error) {
	userIDStr, exits := ctx.Get("user_id")
	if !exits {
		return uuid.UUID{}, errors.New("user_id not found")
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.UUID{}, err
	}
	return userID, nil
}

func (h *Handler) FillWithUsers() {
	h.Repository.FillWithUsers()
}
