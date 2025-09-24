package handler

import (
	"fmt"
	"net/http"
	"strconv"

	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"Backend/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetReactions(ctx *gin.Context) {
	var reactions []ds.Reaction
	var err error

	searchReaction := ctx.Query("reaction_title") // получаем значение из нашего поля
	if searchReaction == "" {                     // если поле поиска пусто, то просто получаем из репозитория все записи
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchReaction) // в ином случае ищем заказ по заголовку
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	resp := make([]apitypes.ReactionJSON, 0, len(reactions))
	for _, r := range reactions {
		resp = append(resp, apitypes.ReactionToJSON(r))
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /reaction/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reaction, err := h.Repository.GetReaction(id)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ReactionToJSON(reaction))
}

func (h *Handler) CreateReaction(ctx *gin.Context) {
	var reactionJSON apitypes.ReactionJSON
	if err := ctx.BindJSON(&reactionJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	reaction, err := h.Repository.CreateReaction(reactionJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/reactions/%v", reaction.ID))
	ctx.JSON(http.StatusCreated, apitypes.ReactionToJSON(reaction))
}

func (h *Handler) ChangeReaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var reactionJSON apitypes.ReactionJSON
	if err := ctx.BindJSON(&reactionJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reaction, err := h.Repository.ChangeReaction(id, reactionJSON)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ReactionToJSON(reaction))
}

func (h *Handler) DeleteReaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteReaction(id)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *Handler) AddReactionToCalculation(ctx *gin.Context) {
	calculation, created, err := h.Repository.GetMassCalculationDraft(h.Repository.GetUserID())
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddReactionToCalculation(calculation.ID, reactionID)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	status := http.StatusOK
	if created {
		ctx.Header("Location", fmt.Sprintf("/calculations/%v", calculation.ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(calculation)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, apitypes.MassCalculationToJSON(calculation, creatorLogin, moderatorLogin))
}

func (h *Handler) UploadImage(ctx *gin.Context) {
	reactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reaction, err := h.Repository.UploadImage(ctx, reactionID, file)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "uploaded",
		"reaction": apitypes.ReactionToJSON(reaction),
	})
}
