package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMassCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// h.errorHandler(ctx, http.StatusBadRequest, err)
		ctx.Redirect(http.StatusFound, "/reactions")
		return
	}

	reactions, calculation, err := h.Repository.GetMassCalculationReactions(id)
	if err != nil {
		// ctx.HTML(http.StatusNotFound, "mass_calculation.html", gin.H{
		// 	"reactions":   nil,
		// 	"calculation": nil,
		// })
		ctx.Redirect(http.StatusFound, "/reactions")
		return
	}

	ctx.HTML(http.StatusOK, "mass_calculation.html", gin.H{
		"reactions":   reactions,
		"calculation": calculation,
	})
}

func (h *Handler) DeleteMassCalculation(ctx *gin.Context) {
	calculationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		// h.errorHandler(ctx, http.StatusBadRequest, err)
		ctx.Redirect(http.StatusFound, "/reactions")
		return
	}

	err = h.Repository.DeleteMassCalculation(calculationID)
	if err != nil {
		// h.errorHandler(ctx, http.StatusInternalServerError, err)
		ctx.Redirect(http.StatusFound, "/reactions")
		return
	}

	ctx.Redirect(http.StatusFound, "/reactions")
}
