package handler

import (
	"FinCoach/internal/app/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"strconv"
)

// /Users/drakenchef/go/bin/swag init -g cmd/main/main.go

// PlanetsList godoc
// @Summary Список планет
// @Description Получение планет и фильтрация при поиске
// @Tags Планеты
// @Produce json
// @Param search query string false "Получаем определённую планету "
// @Success 200 {object} models.PlanetsListResp
// @Failure 400 {object} errorResp "Неверный запрос"
// @Failure 500 {object} errorResp "Внутренняя ошибка сервера"
// @Router /Planets [get]
func (h *Handler) PlanetsList(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	fmt.Println(userID)
	if !exists {
		// Обработка ситуации, когда userid отсутствует в контексте
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("user_id not found in context"))
		return
	}
	// Приведение типа, если необходимо
	var userIDUint uint
	switch v := userID.(type) {
	case uint:
		userIDUint = v
	case int:
		userIDUint = uint(v)
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errors.New("failed to convert user_id to uint"))
			return
		}
		userIDUint = uint(i)
	default:
		h.errorHandler(ctx, http.StatusInternalServerError, errors.New("user_id is not of a supported type"))
		return
	}

	searchQuery := ctx.Query("search")
	if searchQuery == "" {
		planets, err := h.Repository.PlanetsList()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Получаем id заявки пользователя
		userRequestID, err := h.Repository.GetUserRequestID(int(userIDUint))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Planets":   planets,
			"Flight_id": userRequestID,
		})
	} else {
		filteredPlanets, err := h.Repository.SearchPlanet(searchQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Получаем id заявки пользователя
		userRequestID, err := h.Repository.GetUserRequestID(int(userIDUint))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"Planets":   filteredPlanets,
			"Flight_id": userRequestID,
		})
	}
}

func (h *Handler) PlanetById(ctx *gin.Context) {
	id := ctx.Param("id")
	idint, _ := strconv.Atoi(id)
	planets, err := h.Repository.PlanetById(idint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"Planets": planets,
	})
}

// DeletePlanet godoc
// @Summary Удаление планеты
// @Description Удаление планеты по её идентификатору.
// @Security ApiKeyAuth
// @Tags Планеты
// @Accept json
// @Produce json
// @Param request body models.DeletePlanetReq true "ID планеты для удаления"
// @Success 200 {object} models.DeletePlanetRes "Планеты успешно удалена"
// @Failure 400 {object} errorResp "Неверный запрос"
// @Failure 500 {object} errorResp "Внутренняя ошибка сервера"
// @Router /Planets [delete]
func (h *Handler) DeletePlanet(ctx *gin.Context) {
	var request struct {
		ID uint `json:"id"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if request.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if err := h.Repository.DeletePlanet(request.ID); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	h.successHandler(ctx, "Planet_id", request.ID)
}

func (h *Handler) createSpectrum(Planet *models.Planet) (int, error) {
	if Planet.ID != 0 {
		return http.StatusBadRequest, idMustBeEmpty
	}
	if Planet.Description == "" {
		return http.StatusBadRequest, errors.New("planet cannot be empty")
	}
	if err := h.Repository.AddPlanet(Planet); err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}

// UpdatePlanet godoc
// @Summary Обновление информации о планете
// @Security ApiKeyAuth
// @Tags Планеты
// @Description Обновление информации о планете
// @Accept json
// @Produce json
// @Param updated_planet body models.UpdatePlanetReq true "Обновленная информация о планете"
// @Success 200 {object} models.UpdatePlanetResp
// @Failure 400 {object} errorResp "Неверный запрос"
// @Failure 500 {object} errorResp "Внутренняя ошибка сервера"
// @Router /Planets [put]
func (h *Handler) UpdatePlanet(ctx *gin.Context) {
	var updatedPlanet models.Planet
	if err := ctx.BindJSON(&updatedPlanet); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	if updatedPlanet.Image != "" {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New(`image must be empty`))
		return
	}

	if updatedPlanet.ID == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if err := h.Repository.UpdatePlanet(&updatedPlanet); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	h.successHandler(ctx, "updated_Planet", gin.H{
		"id":          updatedPlanet.ID,
		"name":        updatedPlanet.Name,
		"is_delete":   updatedPlanet.IsDelete,
		"description": updatedPlanet.Description,
		"image":       updatedPlanet.Image,
	})
}

// asd=
func (h *Handler) updatePlanet(Planet *models.Planet) (int, error) {
	if Planet.ID == 0 {
		return http.StatusBadRequest, idNotFound
	}
	if err := h.Repository.UpdatePlanet(Planet); err != nil {
		return http.StatusBadRequest, err
	}
	return 0, nil
}

func (h *Handler) AddImage(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	planetID := ctx.Request.FormValue("id")

	if planetID == "" {
		h.errorHandler(ctx, http.StatusBadRequest, idNotFound)
		return
	}
	if header == nil || header.Size == 0 {
		h.errorHandler(ctx, http.StatusBadRequest, headerNotFound)
		return
	}
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	defer func(file multipart.File) {
		errLol := file.Close()
		if errLol != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, errLol)
			return
		}
	}(file)

	h.successAddHandler(ctx, "image", "created image")
}
