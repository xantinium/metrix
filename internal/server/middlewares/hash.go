package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xantinium/metrix/internal/logger"
	"github.com/xantinium/metrix/internal/tools"
)

// HashCheckMiddleware мидлварь для проверки данных
// при помощи хеширования через SHA-256.
func HashCheckMiddleware(privateKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hashedReq := ctx.GetHeader(tools.HashSHA256)
		if hashedReq == "" {
			logger.Errorf("header %q is required", tools.HashSHA256)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		reqBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			logger.Errorf("failed to read request bytes: %v", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var targetHashedReq string
		targetHashedReq, err = tools.CalcSHA256(reqBytes, privateKey)
		if err != nil {
			logger.Errorf("failed to check request by hash: %v", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if hashedReq != targetHashedReq {
			logger.Errorf("hash of request doesn't match to it's content")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// После вызова io.ReadAll требуется восстановить буфер.
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
		ctx.Next()
	}
}

// ResponseHasherMiddleware мидлварь для вычисления
// хеша через SHA-256 и его записи в заголовки ответа.
func ResponseHasherMiddleware(privateKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if privateKey == "" {
			ctx.Next()
			return
		}

		rhw := newResponseHasherWriter(ctx.Writer)
		ctx.Writer = rhw

		ctx.Next()

		hashedRes, err := tools.CalcSHA256(rhw.body.Bytes(), privateKey)
		if err != nil {
			logger.Errorf("failed to calc hash of response: %v", err)
			return
		}

		ctx.Header(tools.HashSHA256, hashedRes)
	}
}

func newResponseHasherWriter(w gin.ResponseWriter) *responseHasherWriter {
	return &responseHasherWriter{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}
}

type responseHasherWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseHasherWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
