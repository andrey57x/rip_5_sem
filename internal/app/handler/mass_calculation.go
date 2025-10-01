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
