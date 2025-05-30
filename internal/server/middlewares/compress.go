package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xantinium/metrix/internal/tools"
)

// CompressMiddleware мидлварь для сжатия данных.
func CompressMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isGZIPSupported(ctx) && isSupportedMIMEType(ctx) {
			// Меняем оригинальный gin.ResponseWriter на новый с поддержкой сжатия.
			cw := newCompressWriter(ctx.Writer)
			ctx.Writer = cw
			defer cw.Close()
		}

		if isRequestCompressed(ctx) {
			// Оборачиваем тело запроса в io.Reader с поддержкой декомпрессии.
			cr, err := newCompressReader(ctx.Request.Body)
			if err != nil {
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
			} else {
				// Меняем тело запроса на новое.
				ctx.Request.Body = cr
				defer cr.Close()
			}
		}

		ctx.Next()
	}
}

const (
	acceptHeader          = "Accept"
	acceptEncodingHeader  = "Accept-Encoding"
	contentTypeHeader     = "Content-Type"
	contentEncodingHeader = "Content-Encoding"
)

// isGZIPSupported проверяет поддержку клиентом сжатия в формате gzip.
func isGZIPSupported(ctx *gin.Context) bool {
	h := ctx.GetHeader(acceptEncodingHeader)

	return h != "" && strings.Contains(h, "gzip")
}

// isRequestCompressed проверяет наличие сжатия запроса в формате gzip.
func isRequestCompressed(ctx *gin.Context) bool {
	h := ctx.GetHeader(contentEncodingHeader)

	return h != "" && strings.Contains(h, "gzip")
}

var supportedMIMETypes = []string{
	"application/json",
	"text/html",
}

// isSupportedMIMEType проверяет заголовоки Accept
// и Content-Type, т.к. не все типы подлежат сжатию.
func isSupportedMIMEType(ctx *gin.Context) bool {
	supported := slices.ContainsFunc(ctx.Request.Header.Values(acceptHeader), func(acceptType string) bool {
		for _, mimeType := range supportedMIMETypes {
			if strings.Contains(acceptType, mimeType) {
				return true
			}
		}

		return false
	})
	if supported {
		return true
	}

	return slices.ContainsFunc(ctx.Request.Header.Values(contentTypeHeader), func(contentTypeType string) bool {
		for _, mimeType := range supportedMIMETypes {
			if strings.Contains(contentTypeType, mimeType) {
				return true
			}
		}

		return false
	})
}

// compressWriter реализует интерфейс gin.ResponseWriter.
type compressWriter struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w gin.ResponseWriter) *compressWriter {
	return &compressWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func (w *compressWriter) Write(p []byte) (int, error) {
	return w.zw.Write(p)
}

func (w *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		w.ResponseWriter.Header().Set(tools.ContentEncoding, "gzip")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *compressWriter) Close() error {
	return w.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser.
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (r compressReader) Read(p []byte) (n int, err error) {
	return r.zr.Read(p)
}

func (r *compressReader) Close() error {
	if err := r.r.Close(); err != nil {
		return err
	}
	return r.zr.Close()
}
