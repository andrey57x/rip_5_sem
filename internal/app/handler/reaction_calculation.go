package handler

import (
	apitypes "Backend/internal/app/api_types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteReactionFromCalculation godoc
// @Summary Delete reaction from calculation
// @Description Удалить реакцию из расчета
// @Tags reactions calculations
// @Accept json
// @Produce json
// @Param calculation_id path int true "Calculation ID"
// @Param reaction_id path int true "Reaction ID"
// @Success 200 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reaction-calculations/{calculation_id}/{reaction_id} [delete]
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
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(calculation)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.CalculationToJSON(calculation, creatorLogin, moderatorLogin))
}

// ChangeReactionCalculation godoc
// @Summary Change reaction calculation
// @Description Изменить поле для реакции в расчете
// @Tags reactions calculations
// @Accept json
// @Produce json
// @Param calculation_id path int true "Calculation ID"
// @Param reaction_id path int true "Reaction ID"
// @Param reaction_calculation body apitypes.ReactionCalculationJSON true "Change reaction calculation"
// @Success 200 {object} apitypes.ReactionCalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reaction-calculations/{calculation_id}/{reaction_id} [put]
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
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ReactionCalculationToJSON(reactionCalculation))

}
