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
	router.GET("/reactions", h.GetReactions)
	router.GET("/reactions/:id", h.GetReaction)
	router.POST("/reactions", h.CreateReaction)
	router.PUT("/reactions/:id", h.ChangeReaction)
	router.DELETE("/reactions/:id", h.DeleteReaction)
	router.POST("/reactions/:id/add-to-calculation", h.AddReactionToCalculation)
	router.POST("/reactions/:id/image", h.UploadImage)

	router.GET("/calculations/calculation-cart", h.GetCalculationCart)
	router.GET("/calculations", h.GetCalculations)
	router.GET("/calculations/:id", h.GetCalculation)
	router.PUT("/calculations/:id", h.ChangeCalculation)
	router.PUT("/calculations/:id/form", h.FormCalculation)
	router.PUT("/calculations/:id/moderate", h.ModerateCalculation)
	router.DELETE("/calculations/:id", h.DeleteCalculation)

	router.DELETE("/reaction-calculations/:calculation_id/:reaction_id", h.DeleteReactionFromCalculation)
	router.PUT("/reaction-calculations/:calculation_id/:reaction_id", h.ChangeReactionCalculation)

	router.POST("/users/sign-up", h.CreateUser)
	router.GET("/users/:id/profile", h.GetProfile)
	router.PUT("/users/:id/profile", h.ChangeProfile)
	router.POST("/users/sign-in", h.SignIn)
	router.POST("/users/sign-out", h.SignOut)
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
