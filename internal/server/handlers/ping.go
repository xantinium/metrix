package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/server/interfaces"
)

// PingHandler реализация хендлера для проверки соединения с БД.
func PingHandler(ctx *gin.Context, s interfaces.Server) (int, string, error) {
	err := s.GetMetricsRepo().CheckDatabase(ctx)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	return http.StatusOK, "", nil
}
