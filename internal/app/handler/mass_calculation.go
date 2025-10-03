package handler

import (
	apitypes "Backend/internal/app/api_types"
	"Backend/internal/app/ds"
	"Backend/internal/app/repository"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MassCalculationResponse struct {
	Reactions       []ReactionWithOutput         `json:"reactions"`
	MassCalculation apitypes.MassCalculationJSON `json:"calculation"`
}

type ReactionWithOutput struct {
	Reaction   apitypes.ReactionJSON `json:"reaction"`
	OutputMass float32               `json:"output_mass"`
	InputMass  float32               `json:"input_mass"`
}

// GetMassCalculation
// @Summary Получить расчет по ID
// @Description Возвращает расчет и список реакций, входящих в неё, с рассчитанными массами.
// @Tags mass-calculations
// @Produce json
// @Param id path int true "ID расчета"
// @Success 200 {object} MassCalculationResponse "Объект расчета с реакциями"
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 403 {object} map[string]string "Доступ запрещён"
// @Failure 404 {object} map[string]string "Расчет не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations/{id} [get]
func (h *Handler) GetMassCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactions, calculation, err := h.Repository.GetMassCalculationReactions(id)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err == repository.ErrorDeleted {
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

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(calculation)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	calculationJSON := apitypes.MassCalculationToJSON(calculation, creatorLogin, moderatorLogin)

	reactionsWithOutput := make([]ReactionWithOutput, len(reactions))

	for i, reaction := range reactions {
		output, err := h.Repository.GetReactionCalculation(reaction.ID, calculation.ID)

		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}

		outputMass := output.OutputMass
		inputMass := output.InputMass

		reactionsWithOutput[i] = ReactionWithOutput{
			Reaction:   apitypes.ReactionToJSON(reaction),
			OutputMass: outputMass,
			InputMass:  inputMass,
		}
	}

	massCalculationResponse := MassCalculationResponse{
		Reactions:       reactionsWithOutput,
		MassCalculation: calculationJSON,
	}

	ctx.JSON(http.StatusOK, massCalculationResponse)
}

// GetMassCalculationCart
// @Summary Получить карточку черновика расчета для пользователя
// @Description Возвращает количество реакций в текущей корзине пользователя и ID черновика расчета (если есть).
// @Tags mass-calculations
// @Produce json
// @Success 200 {object} map[string]interface{} "Поля: id (int, -1 если нет черновика), reactions_count (int)"
// @Failure 400 {object} map[string]string "Некорректный запрос (например, неверный токен)"}
// @Failure 404 {object} map[string]string "Черновик не найден"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations/calculation-cart [get]
func (h *Handler) GetMassCalculationCart(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	reactionsCount := h.Repository.GetCartCount(userID)

	if reactionsCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"id":              -1,
			"reactions_count": reactionsCount,
		})
		return
	}

	calculation, err := h.Repository.CheckCurrentMassCalculationDraft(userID)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reactions_count": reactionsCount,
		"id":              calculation.ID,
	})
}

// GetMassCalculations
// @Summary Получить список расчетов
// @Description Возвращает список расчетов с возможностью фильтрации по диапазону дат и статусу.
// @Tags mass-calculations
// @Produce json
// @Param from-date query string false "Нижняя граница даты (YYYY-MM-DD)"
// @Param to-date query string false "Верхняя граница даты (YYYY-MM-DD)"
// @Param status query string false "Статус расчета (draft, formed, moderated, deleted)"
// @Success 200 {array} apitypes.MassCalculationJSON "Список расчетов"
// @Failure 400 {object} map[string]string "Неверный формат даты или параметров запроса"}
// @Failure 404 {object} map[string]string "Не найдены записи"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations [get]
func (h *Handler) GetMassCalculations(ctx *gin.Context) {
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

	calculations, err := h.Repository.GetMassCalculations(from, to, status)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	calculations = h.filterCalculationsByAuth(calculations, ctx)

	resp := make([]apitypes.MassCalculationJSON, 0, len(calculations))
	for _, c := range calculations {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err == repository.ErrorNotFound {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, apitypes.MassCalculationToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}

// ChangeMassCalculation
// @Summary Изменить поля расчета
// @Description Обновляет расчет по ID: принимает JSON расчета и возвращает обновлённый объект.
// @Tags mass-calculations
// @Accept json
// @Produce json
// @Param id path int true "ID расчета"
// @Param calculation body apitypes.MassCalculationJSON true "Тело запроса — объект расчета"
// @Success 200 {object} apitypes.MassCalculationJSON "Обновлённый расчет"
// @Failure 400 {object} map[string]string "Неверный формат запроса или тела"
// @Failure 403 {object} map[string]string "Доступ запрещён"
// @Failure 404 {object} map[string]string "Расчет не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations/{id} [put]
func (h *Handler) ChangeMassCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var calculationJSON apitypes.MassCalculationJSON
	if err := ctx.BindJSON(&calculationJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.ChangeMassCalculation(id, calculationJSON)
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

// FormMassCalculation
// @Summary Сформировать расчет
// @Description Переводит расчет в статус formed — доступен только создателю.
// @Tags mass-calculations
// @Produce json
// @Param id path int true "ID расчета"
// @Success 200 {object} apitypes.MassCalculationJSON "Расчет успешно сформирована"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 403 {object} map[string]string "Только создатель может формировать расчет"
// @Failure 404 {object} map[string]string "Расчет не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations/{id}/form [put]
func (h *Handler) FormMassCalculation(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.GetSingleMassCalculation(id)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if calculation.CreatorID != userID {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("only creator can form mass calculation"))
		return
	}

	status := "formed"

	calculation, err = h.Repository.FormMassCalculation(id, status)

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

// DeleteMassCalculation
// @Summary Удалить расчет (логическое удаление)
// @Description Помечает расчет как удалённый. Доступ — владелец или модератор. Возвращает {"message":"Calculation deleted"}.
// @Tags mass-calculations
// @Produce json
// @Param id path int true "ID расчета"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 403 {object} map[string]string "Доступ запрещён"
// @Failure 404 {object} map[string]string "Расчет не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations/{id} [delete]
func (h *Handler) DeleteMassCalculation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	calculation, err := h.Repository.GetSingleMassCalculation(id)
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

	status := "deleted"

	_, err = h.Repository.FormMassCalculation(id, status)

	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Calculation deleted"})
}

// ModerateMassCalculation
// @Summary Модерировать расчет (только модератор)
// @Description Устанавливает статус расчета — только модератор может это сделать.
// @Tags mass-calculations
// @Accept json
// @Produce json
// @Param id path int true "ID расчета"
// @Param status body apitypes.StatusJSON true "Тело запроса с полем status"
// @Success 200 {object} apitypes.MassCalculationJSON "Расчет после модерации"
// @Failure 400 {object} map[string]string "Некорректные входные данные"}
// @Failure 401 {object} map[string]string "Неавторизован"}
// @Failure 403 {object} map[string]string "Только модератор может модерацию"
// @Failure 404 {object} map[string]string "Расчет не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /mass-calculations/{id}/moderate [put]
func (h *Handler) ModerateMassCalculation(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var statusJSON apitypes.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	if !user.IsModerator {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}

	calculation, err := h.Repository.ModerateMassCalculation(id, statusJSON.Status, userID)

	if err == repository.ErrorNotFound {
		h.errorHandler(ctx, http.StatusNotFound, err)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
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

func (h *Handler) filterCalculationsByAuth(calculations []ds.MassCalculation, ctx *gin.Context) []ds.MassCalculation {
	userID, err := getUserID(ctx)
	if err != nil {
		return []ds.MassCalculation{}
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrorNotFound {
		return []ds.MassCalculation{}
	}
	if err != nil {
		return []ds.MassCalculation{}
	}

	if user.IsModerator {
		return calculations
	}

	for _, calculation := range calculations {
		if calculation.CreatorID == userID {
			return []ds.MassCalculation{calculation}
		}
	}
	return []ds.MassCalculation{}
}

func (h *Handler) hasAccessToCalculation(creatorID uuid.UUID, ctx *gin.Context) bool {
	userID, err := getUserID(ctx)
	if err != nil {
		return false
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrorNotFound {
		return false
	}
	if err != nil {
		return false
	}

	return creatorID == userID || user.IsModerator
}
