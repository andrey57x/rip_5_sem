package handler

import (
	apitypes "Backend/internal/app/api_types"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) DeleteCalculation(ctx *gin.Context) {
	calculationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteCalculation(calculationID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/reactions")
}

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

func (h *Handler) FormCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	var status string
	if ctx.Request.Method == "PUT" {
		status = "formed"
	} else {
		status = "deleted"
	}

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
