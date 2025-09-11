package handler

import (
	"net/http"
	"strconv"

	"Backend/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetReactions(ctx *gin.Context) {
	var reactions []ds.Reaction
	var err error

	searchReaction := ctx.Query("reaction_title") // получаем значение из нашего поля
	if searchReaction == "" {                     // если поле поиска пусто, то просто получаем из репозитория все записи
		reactions, err = h.Repository.GetReactions()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		reactions, err = h.Repository.GetReactionsByTitle(searchReaction) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	calculation, _ := h.Repository.CheckCurrentCalculationDraft(h.Repository.GetUser())

	ctx.HTML(http.StatusOK, "reactions.html", gin.H{
		"reactions":                reactions,
		"reaction_title":           searchReaction,
		"reactions_in_calculation": h.Repository.GetCartCount(),
		"calculation_id":           calculation.ID,
	})
}

func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /reaction/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	reaction, err := h.Repository.GetReaction(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "single_reaction.html", gin.H{
		"reaction": reaction,
	})
}

func (h *Handler) AddReactionToCalculation(ctx *gin.Context) {
	calculation, err := h.Repository.GetCalculationDraft(h.Repository.GetUser())
	calculationID := calculation.ID
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddReactionToCalculation(calculationID, reactionID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/reactions")
}
