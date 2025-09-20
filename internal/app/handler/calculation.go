package handler

import (
	apitypes "Backend/internal/app/api_types"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetCalculation godoc
// @Summary Get calculation
// @Description Получить расчет по id
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path int true "Calculation ID"
// @Success 200 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations/{id} [get]
func (h *Handler) GetCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	reactions, calculation, err := h.Repository.GetCalculationReactions(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	resp := make([]apitypes.ReactionJSON, 0, len(reactions))
	for _, r := range reactions {
		resp = append(resp, apitypes.ReactionToJSON(r))
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(calculation)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"calculation": apitypes.CalculationToJSON(calculation, creatorLogin, moderatorLogin),
		"reactions":   resp,
	})
}

// GetCalculationCart godoc
// @Summary Get calculation cart
// @Description Получить черновик расчета
// @Tags calculations
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations/calculation-cart [get]
func (h *Handler) GetCalculationCart(ctx *gin.Context) {
	reactionsCount := h.Repository.GetCartCount(h.Repository.GetUserID())

	if reactionsCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":          "no_draft",
			"reactions_count": reactionsCount,
		})
		return
	}

	calculation, err := h.Repository.CheckCurrentCalculationDraft(h.Repository.GetUserID())
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":          "draft",
		"reactions_count": reactionsCount,
		"id":              calculation.ID,
	})
}

// GetCalculations godoc
// @Summary Get calculations
// @Description Получить расчеты
// @Tags calculations
// @Accept json
// @Produce json
// @Param from-date query string false "From date"
// @Param to-date query string false "To date"
// @Param status query string false "Status"
// @Success 200 {object} []apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations [get]
func (h *Handler) GetCalculations(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}

	toDate := ctx.Query("to-date")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	calculations, err := h.Repository.GetCalculations(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	resp := make([]apitypes.CalculationJSON, 0, len(calculations))
	for _, c := range calculations {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, apitypes.CalculationToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}


// ChangeCalculation godoc
// @Summary Change calculation
// @Description Изменить расчет
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path int true "Calculation ID"
// @Param calculation body apitypes.CalculationJSON true "Change calculation"
// @Success 200 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations/{id} [put]
func (h *Handler) ChangeCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var calculationJSON apitypes.CalculationJSON
	if err := ctx.BindJSON(&calculationJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.ChangeCalculation(id, calculationJSON)
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

// FormCalculation godoc
// @Summary Form calculation
// @Description Сформировать расчет
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path int true "Calculation ID"
// @Success 200 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations/{id}/form [put]
func (h *Handler) FormCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	status := "formed"

	calculation, err := h.Repository.FormCalculation(id, status)

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

// DeleteCalculation godoc
// @Summary Delete calculation
// @Description Удалить расчет
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path int true "Calculation ID"
// @Success 200 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations/{id} [delete]
func (h *Handler) DeleteCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	status := "deleted"

	_, err = h.Repository.FormCalculation(id, status)

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Calculation deleted"})
}

// ModerateCalculation godoc
// @Summary Moderate calculation
// @Description Модерировать расчет
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path int true "Calculation ID"
// @Param status body apitypes.StatusJSON true "Moderate calculation"
// @Success 200 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /calculations/{id}/moderate [put]
func (h *Handler) ModerateCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var statusJSON apitypes.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.ModerateCalculation(id, statusJSON.Status)

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
