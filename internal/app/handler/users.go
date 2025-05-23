package handler

import (
	"FinCoach/internal/app/models"
	"FinCoach/internal/app/role"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) Login(ctx *gin.Context) {
	cfg := h.Config
	req := &models.LoginReq{}

	if err := json.NewDecoder(ctx.Request.Body).Decode(req); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		h.errorHandler(ctx, http.StatusForbidden, err)
		return
	}

	fmt.Println(req.Login, user.Login, user.Password, generateHashString(req.Password))
	if req.Login == user.Login && user.Password == generateHashString(req.Password) {
		token := jwt.NewWithClaims(cfg.JWT.SigningMethod, &models.JWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(cfg.JWT.ExpiresIn).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "bitop-admin",
			},
			UserID: user.ID,
			Role:   user.Role,
		})

		if token == nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errors.New("token Is Nil"))
			return
		}

		strToken, err := token.SignedString([]byte(cfg.JWT.Token))
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errors.New("cannot create token"))
			return
		}

		ctx.JSON(http.StatusOK, models.LoginResp{
			ExpiresIn:   cfg.JWT.ExpiresIn,
			AccessToken: strToken,
			TokenType:   "Bearer",
			Role:        user.Role,
			Username:    user.UserName,
		})
	}

	ctx.AbortWithStatus(http.StatusForbidden)
}

func (h *Handler) UsersList(ctx *gin.Context) {
	users, err := h.Repository.UsersList()
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "users", users)
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param request body models.RegisterReq true "Детали регистрации"
// @Router /signup [post]
func (h *Handler) Register(ctx *gin.Context) {
	type registerReq struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Username string `json:"user_name"`
	}

	type registerResp struct {
		Ok bool `json:"ok"`
	}

	req := &registerReq{}

	err := json.NewDecoder(ctx.Request.Body).Decode(req)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if req.Password == "" {
		//h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("pass is empty"))
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("pass is empty"))
		return
	}

	if req.Login == "" {
		//h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("login is empty"))
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("login is empty"))
		return
	}
	if req.Username == "" {
		//h.errorHandler(ctx, http.StatusBadRequest, fmt.Errorf("username is empty"))
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("username is empty"))
		return
	}

	err = h.Repository.Register(&models.Users{
		Role:     role.Buyer,
		Login:    req.Login,
		UserName: req.Username,
		Password: generateHashString(req.Password),
	})

	if err != nil {
		// Проверяем, является ли ошибка нарушением уникального ограничения
		if strings.Contains(err.Error(), "duplicate key value") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": "Пользователь с таким логином уже существует",
			})
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &registerResp{
		Ok: true,
	})
}

// Logout godoc
// @Summary Выход пользователя
// @Description Завершение сеанса текущего пользователя.
// @Tags Пользователи
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {string} string "Успешный выход"
// @Failure 400 {object} errorResp "Неверный запрос"
// @Failure 401 {object} errorResp "Неверные учетные данные"
// @Failure 500 {object} errorResp "Внутренняя ошибка сервера"
// @Router /logout [get]
func (h *Handler) Logout(ctx *gin.Context) {
	jwtStr := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Redis.WriteJWTToBlacklist(ctx.Request.Context(), jwtStr, h.Config.JWT.ExpiresIn)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// MARK: - Inner functions

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
