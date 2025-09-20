package handler

import (
	"fmt"
	"net/http"
	"strconv"

	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// GetReactions godoc
// @Summary List reactions
// @Description Получить список реакций, можно фильтровать по title
// @Tags reactions
// @Accept json
// @Produce json
// @Param reaction_title query string false "Search by title"
// @Success 200 {array} apitypes.ReactionJSON
// @Failure 500 {object} map[string]string
// @Router /reactions [get]
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

// GetReaction godoc
// @Summary Get reaction
// @Description Получить реакцию по id
// @Tags reactions
// @Accept json
// @Produce json
// @Param id path int true "Reaction ID"
// @Success 200 {object} apitypes.ReactionJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reactions/{id} [get]
func (h *Handler) GetReaction(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /reaction/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reaction, err := h.Repository.GetReaction(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ReactionToJSON(reaction))
}

// CreateReaction godoc
// @Summary Create reaction
// @Description Создать реакцию
// @Tags reactions
// @Accept json
// @Produce json
// @Param reaction body apitypes.ReactionJSON true "Create reaction"
// @Success 201 {object} apitypes.ReactionJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reactions [post]
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

// ChangeReaction godoc
// @Summary Change reaction
// @Description Изменить реакцию
// @Tags reactions
// @Accept json
// @Produce json
// @Param id path int true "Reaction ID"
// @Param reaction body apitypes.ReactionJSON true "Change reaction"
// @Success 200 {object} apitypes.ReactionJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reactions/{id} [put]
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
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ReactionToJSON(reaction))
}

// DeleteReaction godoc
// @Summary Delete reaction
// @Description Удалить реакцию
// @Tags reactions
// @Accept json
// @Produce json
// @Param id path int true "Reaction ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reactions/{id} [delete]
func (h *Handler) DeleteReaction(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeleteReaction(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

// AddReactionToCalculation godoc
// @Summary Add reaction to calculation
// @Description Добавить реакцию в калькуляцию
// @Tags calculations
// @Accept json
// @Produce json
// @Param id path int true "Reaction ID"
// @Success 200 {object} apitypes.CalculationJSON
// @Success 201 {object} apitypes.CalculationJSON
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reactions/{id}/add-to-calculation [post]
func (h *Handler) AddReactionToCalculation(ctx *gin.Context) {
	calculation, created, err := h.Repository.GetCalculationDraft(h.Repository.GetUserID())
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
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, apitypes.CalculationToJSON(calculation, creatorLogin, moderatorLogin))
}

// UploadImage godoc
// @Summary Upload image
// @Description Загрузить изображение
// @Tags reactions
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Reaction ID"
// @Param image formData file true "Image"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reactions/{id}/image [post]
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
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
		"reaction": apitypes.ReactionToJSON(reaction),
	})
}
