package tools

import (
	"net/http"

	"github.com/xantinium/metrix/internal/logger"
)

// RunProfilingServer запускает веб-сервер с интерактивной
// панелью профилирования.
func RunProfilingServer() {
	go func() {
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			logger.Errorf("failed to start pprof server: %v", err)
		}
	}()
}
