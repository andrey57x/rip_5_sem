package handler

import (
	"Backend/internal/app/repository"
	"net/http"

	_ "Backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")
	unauthorized.POST("/users/sign-in", h.SignIn)
	unauthorized.POST("/users/sign-up", h.SignUp)
	unauthorized.GET("/reactions", h.GetReactions)
	unauthorized.GET("/reactions/:id", h.GetReaction)
	unauthorized.GET("/mass-calculations/mass-calculation-cart-icon", h.OptionalAuthMiddleware(), h.GetIconCart)
	unauthorized.POST("/mass-calculations/:id/update-result", h.UpdateCalculationResult)

	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))

	authorized.POST("/reactions/:id/add-to-calculation", h.AddReactionToCalculation)

	authorized.GET("/mass-calculations", h.GetMassCalculations)
	authorized.GET("/mass-calculations/:id", h.GetMassCalculation)
	authorized.PUT("/mass-calculations/:id", h.ChangeMassCalculation)
	authorized.PUT("/mass-calculations/:id/form", h.FormMassCalculation)
	authorized.DELETE("/mass-calculations/:id", h.DeleteMassCalculation)

	authorized.DELETE("/reaction-calculations/:mass_calculation_id/:reaction_id", h.DeleteReactionFromCalculation)
	authorized.PUT("/reaction-calculations/:mass_calculation_id/:reaction_id", h.ChangeReactionCalculation)

	authorized.GET("/users/:id/profile", h.GetProfile)
	authorized.PUT("/users/:id/profile", h.ChangeProfile)
	authorized.POST("/users/sign-out", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/mass-calculations/:id/moderate", h.ModerateMassCalculation)
	moderator.POST("/reactions", h.CreateReaction)
	moderator.PUT("/reactions/:id", h.ChangeReaction)
	moderator.DELETE("/reactions/:id", h.DeleteReaction)
	moderator.POST("/reactions/:id/image", h.UploadImage)

	// нужно
	// для
	// swagger
	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static/styles", "./static/styles")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
