package middlewares

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/xantinium/metrix/internal/tools"
)

// DecryptMiddleware мидлварь для дешифроки данных.
// Используется алгоритм AES.
func DecryptMiddleware(privateKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Оборачиваем тело запроса в io.Reader с поддержкой дешифроки.
		dr := newDecryptReader(ctx.Request.Body, privateKey)

		// Меняем тело запроса на новое.
		ctx.Request.Body = dr
		defer dr.Close()

		ctx.Next()
	}
}

// decryptReader реализует интерфейс io.ReadCloser.
type decryptReader struct {
	io.ReadCloser

	privateKey string
}

func newDecryptReader(r io.ReadCloser, privateKey string) *decryptReader {
	return &decryptReader{
		ReadCloser: r,
		privateKey: privateKey,
	}
}

func (r *decryptReader) Read(p []byte) (int, error) {
	var message []byte
	n, err := r.ReadCloser.Read(message)
	if err != nil {
		return n, err
	}

	message, err = tools.Decrypt(r.privateKey, message)
	if err != nil {
		return n, err
	}

	copy(p, message)

	return n, nil
}
