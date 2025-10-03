package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"Backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetReactions
// @Summary Получить список реакций
// @Description Возвращает все реакции или фильтрует по параметру reaction_title.
// @Tags reactions
// @Produce json
// @Param reaction_title query string false "Фильтр по заголовку реакции"
// @Success 200 {array} apitypes.ReactionJSON "Список реакций"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
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

// GetReaction
// @Summary Получить реакцию по ID
// @Description Возвращает реакцию по id.
// @Tags reactions
// @Produce json
// @Param id path int true "ID реакции"
// @Success 200 {object} apitypes.ReactionJSON "Реакция"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Реакция не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
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

// CreateReaction
// @Summary Создать реакцию
// @Description Создаёт новую реакцию, возвращает Location в заголовках и объект реакции.
// @Tags reactions
// @Accept json
// @Produce json
// @Param reaction body apitypes.ReactionJSON true "Данные реакции"
// @Success 201 {object} apitypes.ReactionJSON "Реакция создана"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
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

// ChangeReaction
// @Summary Обновить реакцию
// @Description Обновляет реакцию по ID.
// @Tags reactions
// @Accept json
// @Produce json
// @Param id path int true "ID реакции"
// @Param reaction body apitypes.ReactionJSON true "Новые данные реакции"
// @Success 200 {object} apitypes.ReactionJSON "Обновлённая реакция"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Реакция не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
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

// DeleteReaction
// @Summary Удалить реакцию
// @Description Логическое удаление реакции. Возвращает {"status":"deleted"}.
// @Tags reactions
// @Produce json
// @Param id path int true "ID реакции"
// @Success 200 {object} map[string]string "status"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Реакция не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /reactions/{id} [delete]
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

// AddReactionToCalculation
// @Summary Добавить реакцию в текущий черновик расчета
// @Description Добавляет реакцию в черновик пользователя (если черновика нет — создаёт). Возвращает сам расчет.
// @Tags mass-calculations
// @Produce json
// @Param id path int true "ID реакции"
// @Success 200 {object} apitypes.MassCalculationJSON "Если добавление в существующий черновик"
// @Success 201 {object} apitypes.MassCalculationJSON "Если был создан новый черновик"
// @Failure 400 {object} map[string]string "Ошибка запроса или авторизации"
// @Failure 404 {object} map[string]string "Реакция/калькуляция не найдены"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /reactions/{id}/add-to-calculation [post]
func (h *Handler) AddReactionToCalculation(ctx *gin.Context) {
	userIDStr, exits := ctx.Get("user_id")
	if !exits {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("user_id not found"))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, created, err := h.Repository.GetMassCalculationDraft(userID)
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

// UploadImage
// @Summary Загрузить изображение для реакции
// @Description Загружает файл изображения и возвращает объект вида {"status":"uploaded", "reaction": <ReactionJSON>}.
// @Tags reactions
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID реакции"
// @Param image formData file true "Изображение"
// @Success 200 {object} map[string]interface{} "status (string) и reaction (apitypes.ReactionJSON)"
// @Failure 400 {object} map[string]string "Неверный запрос/файл"
// @Failure 404 {object} map[string]string "Реакция не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
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
