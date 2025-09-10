package handler

import (
	"Backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetReactions(ctx *gin.Context) {
	var reactions []repository.Reaction
	var err error

	searchReaction := ctx.Query("reaction_title")
	if searchReaction == "" {
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchReaction)
		if err != nil {
			logrus.Error(err)
		}
	}

	var calculationsCount int
	var currCalculationId int
	calculationsCount, currCalculationId, err = h.Repository.CurrentCalculation()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "reactions.html", gin.H{
		"reactions":                reactions,
		"reaction_title":           searchReaction,
		"reactions_in_calculation": calculationsCount,
		"mass_calculation_id":      currCalculationId,
	})
}

func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	reaction, err := h.Repository.GetReaction(id)
	if err != nil {
		logrus.Error(err)
	}

	var calculationsCount int
	var currCalculationId int
	calculationsCount, currCalculationId, err = h.Repository.CurrentCalculation()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "single_reaction.html", gin.H{
		"reaction":                 reaction,
		"reactions_in_calculation": calculationsCount,
		"mass_calculation_id":      currCalculationId,
	})
}

func (h *Handler) GetMassCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	reactions, err := h.Repository.GetCalculationReactions(id)
	if err != nil {
		logrus.Error(err)
	}

	var calculationsCount int
	var currCalculationId int
	calculationsCount, currCalculationId, err = h.Repository.CurrentCalculation()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "mass_calculation.html", gin.H{
		"reactions":                reactions,
		"reactions_in_calculation": calculationsCount,
		"mass_calculation_id":      currCalculationId,
	})
}
