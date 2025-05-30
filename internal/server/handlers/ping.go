package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/server/interfaces"
)

// PingHandler реализация хендлера для проверки соединения с БД.
// @Tags Database
// @Summary Запрос на проверку соединения с БД.
// @Description Запрос на проверку соединения с БД.
// @ID ping
// @Produce text/plain
// @Success 200 {string} string
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /ping [get]
func PingHandler(ctx *gin.Context, s interfaces.Server) (int, string, error) {
	err := s.GetMetricsRepo().CheckDatabase(ctx)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, "", nil
}
