package handler

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteReactionFromCalculation
// @Summary Удалить реакцию из расчета
// @Description Удаляет связь реакции и расчета и возвращает обновлённый расчет.
// @Tags reaction-calculations
// @Produce json
// @Param mass_calculation_id path int true "ID расчета"
// @Param reaction_id path int true "ID реакции"
// @Success 200 {object} apitypes.MassCalculationJSON "Обновлённый расчет"
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 403 {object} map[string]string "Доступ запрещён"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /reaction-calculations/{mass_calculation_id}/{reaction_id} [delete]
func (h *Handler) DeleteReactionFromCalculation(ctx *gin.Context) {
	calculationID, err := strconv.Atoi(ctx.Param("mass_calculation_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionID, err := strconv.Atoi(ctx.Param("reaction_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.GetSingleMassCalculation(calculationID)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !h.hasAccessToCalculation(calculation.CreatorID, ctx) {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}

	calculation, err = h.Repository.DeleteReactionFromCalculation(calculationID, reactionID)
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

// ChangeReactionCalculation
// @Summary Изменить поле связи многие ко многим
// @Description Обновляет поле расчёта реакции (output mass).
// @Tags reaction-calculations
// @Accept json
// @Produce json
// @Param mass_calculation_id path int true "ID расчета"
// @Param reaction_id path int true "ID реакции"
// @Param body body apitypes.ReactionCalculationJSON true "Новые значения расчёта"
// @Success 200 {object} apitypes.ReactionCalculationJSON "Обновлённый расчёт"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещён"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /reaction-calculations/{mass_calculation_id}/{reaction_id} [put]
func (h *Handler) ChangeReactionCalculation(ctx *gin.Context) {
	calculationID, err := strconv.Atoi(ctx.Param("mass_calculation_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionID, err := strconv.Atoi(ctx.Param("reaction_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.GetSingleMassCalculation(calculationID)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !h.hasAccessToCalculation(calculation.CreatorID, ctx) {
		h.errorHandler(ctx, http.StatusForbidden, err)
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
