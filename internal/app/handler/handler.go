package handler

import (
	"Backend/internal/app/repository"

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

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты, чтобы не писать все в одном месте
func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api/v1")

	api.GET("/reactions", h.GetReactions)
	api.GET("/reactions/:id", h.GetReaction)
	api.POST("/reactions", h.CreateReaction)
	api.PUT("/reactions/:id", h.ChangeReaction)
	api.DELETE("/reactions/:id", h.DeleteReaction)
	api.POST("/reactions/:id/add-to-calculation", h.AddReactionToCalculation)
	api.POST("/reactions/:id/image", h.UploadImage)

	api.GET("/mass-calculations/calculation-cart", h.GetMassCalculationCart)
	api.GET("/mass-calculations", h.GetMassCalculations)
	api.GET("/mass-calculations/:id", h.GetMassCalculation)
	api.PUT("/mass-calculations/:id", h.ChangeMassCalculation)
	api.PUT("/mass-calculations/:id/form", h.FormMassCalculation)
	api.PUT("/mass-calculations/:id/moderate", h.ModerateMassCalculation)
	api.DELETE("/mass-calculations/:id", h.DeleteMassCalculation)

	api.DELETE("/reaction-calculations/:mass_calculation_id/:reaction_id", h.DeleteReactionFromCalculation)
	api.PUT("/reaction-calculations/:mass_calculation_id/:reaction_id", h.ChangeReactionCalculation)

	api.POST("/users/sign-up", h.CreateUser)
	api.GET("/users/:id/profile", h.GetProfile)
	api.PUT("/users/:id/profile", h.ChangeProfile)
	api.POST("/users/sign-in", h.SignIn)
	api.POST("/users/sign-out", h.SignOut)
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static/styles", "./static/styles")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
