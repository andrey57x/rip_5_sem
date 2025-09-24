package handler

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteReactionFromCalculation(ctx *gin.Context) {
	calculationID, err := strconv.Atoi(ctx.Param("calculation_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionID, err := strconv.Atoi(ctx.Param("reaction_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.DeleteReactionFromCalculation(calculationID, reactionID)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
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

	ctx.JSON(http.StatusOK, apitypes.MassCalculationToJSON(calculation, creatorLogin, moderatorLogin))
}

func (h *Handler) ChangeReactionCalculation(ctx *gin.Context) {
	calculationID, err := strconv.Atoi(ctx.Param("calculation_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionID, err := strconv.Atoi(ctx.Param("reaction_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var reactionCalculationJSON apitypes.ReactionCalculationJSON
	if err := ctx.BindJSON(&reactionCalculationJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionCalculation, err := h.Repository.ChangeReactionCalculation(calculationID, reactionID, reactionCalculationJSON)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ReactionCalculationToJSON(reactionCalculation))

}
