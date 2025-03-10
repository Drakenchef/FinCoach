package handler

import (
	"FinCoach/internal/app/config"
	redis2 "FinCoach/internal/app/redis"
	"FinCoach/internal/app/repository"
	"FinCoach/internal/app/role"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
	Logger     *logrus.Logger
	Repository *repository.Repository
	Config     *config.Config
	Redis      *redis2.Client
}

func NewHandler(l *logrus.Logger, r *repository.Repository, conf *config.Config, red *redis2.Client) *Handler {
	return &Handler{
		Logger:     l,
		Repository: r,
		Config:     conf,
		Redis:      red,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	h.UserCRUD(router)
	h.CreditCRUD(router)
	h.SpendingCRUD(router)
	h.Balance(router)
	h.GoalCRUD(router)
	h.MainPage(router)
	h.CategoriesCRUD(router)
	registerStatic(router)
}
func (h *Handler) UserCRUD(router *gin.Engine) {
	router.GET("/users", h.UsersList)
	router.POST("/login", h.Login)
	router.POST("/signup", h.Register)
	router.GET("/logout", h.Logout)
}

func (h *Handler) CreditCRUD(router *gin.Engine) {
	router.POST("/AddCredit", h.WithIdCheck(role.Buyer, role.Moder), h.GoalCheckMiddleware(), h.AddCredit)
	router.GET("/Credits", h.WithIdCheck(role.Buyer, role.Moder), h.GetCredits)
	router.GET("/Credit/:id", h.WithIdCheck(role.Buyer, role.Moder), h.GetCreditByID)
	router.PUT("/Credit/:id", h.WithIdCheck(role.Buyer, role.Moder), h.UpdateCreditByID)
	router.DELETE("/Credit/:id", h.WithIdCheck(role.Buyer, role.Moder), h.DeleteCreditByID)
}

func (h *Handler) SpendingCRUD(router *gin.Engine) {
	router.POST("/AddSpending", h.WithIdCheck(role.Buyer, role.Moder), h.AddSpending)
	router.GET("/Spendings", h.WithIdCheck(role.Buyer, role.Moder), h.GetSpendings)
	router.GET("/Spending/:id", h.WithIdCheck(role.Buyer, role.Moder), h.GetSpendingByID)
	router.PUT("/Spending/:id", h.WithIdCheck(role.Buyer, role.Moder), h.UpdateSpendingByID)
	router.DELETE("/Spending/:id", h.WithIdCheck(role.Buyer, role.Moder), h.DeleteSpendingByID)
}

func (h *Handler) Balance(router *gin.Engine) {
	router.GET("/Balance", h.WithIdCheck(role.Buyer, role.Moder), h.GetBalance)
}

func (h *Handler) GoalCRUD(router *gin.Engine) {
	router.POST("/AddGoal", h.WithIdCheck(role.Buyer, role.Moder), h.AddGoal)
	router.GET("/Goals", h.WithIdCheck(role.Buyer, role.Moder), h.GetGoals)
	router.GET("/Goal/:id", h.WithIdCheck(role.Buyer, role.Moder), h.GetGoalByID)
	router.PUT("/Goal/:id", h.WithIdCheck(role.Buyer, role.Moder), h.UpdateGoalByID)
	router.PUT("/CurrentGoal/:id", h.WithIdCheck(role.Buyer, role.Moder), h.SelectCurrentGoalByID)
	router.GET("/CurrentGoal", h.WithIdCheck(role.Buyer, role.Moder), h.GetCurrentGoal)
	router.DELETE("/Goal/:id", h.WithIdCheck(role.Buyer, role.Moder), h.DeleteGoalByID)
}

func (h *Handler) CategoriesCRUD(router *gin.Engine) {
	router.POST("/AddCategory", h.WithIdCheck(role.Buyer, role.Moder), h.AddCategory)
	router.GET("/Categories", h.WithIdCheck(role.Buyer, role.Moder), h.GetCategories)
	router.GET("/Category/:id", h.WithIdCheck(role.Buyer, role.Moder), h.GetCategoryByID)
	router.PUT("/Category/:id", h.WithIdCheck(role.Buyer, role.Moder), h.UpdateCategoryByID)
	router.DELETE("/Category/:id", h.WithIdCheck(role.Buyer, role.Moder), h.DeleteCategoryByID)
}

func (h *Handler) MainPage(router *gin.Engine) {
	router.GET("/", h.WithIdCheck(role.Buyer, role.Moder), h.GetRecommendation)
}

// request status

func registerStatic(router *gin.Engine) {
	router.Static("/static", "./static")
	router.Static("/img", "./static")
}

// MARK: - Error handler
type errorResp struct {
	Status      string `json:"status" example:"error"`
	Description string `json:"description" example:"Описание ошибки"`
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	h.Logger.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *Handler) successHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		key:      data,
	})
}

func (h *Handler) successAddHandler(ctx *gin.Context, key string, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		key:      data,
	})
}

func (h *Handler) Ping(gCtx *gin.Context) {
	gCtx.String(http.StatusOK, "Hello, my friend!")
}
