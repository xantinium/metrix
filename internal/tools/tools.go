// Пакет tools содержит вспомогательные методы
// и структуры общего назначения.
package tools

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"strconv"
)

const (
	AcceptEncoding  = "Accept-Encoding"
	ContentEncoding = "Content-Encoding"
	ContentType     = "Content-Type"
	HashSHA256      = "HashSHA256"
)

// FloatToStr конвертирует float64 в строку.
func FloatToStr(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// IntToStr конвертирует int64 в строку.
func IntToStr(v int64) string {
	return strconv.FormatInt(v, 10)
}

// StrToFloat конвертирует строку в float64.
func StrToFloat(v string) (float64, error) {
	return strconv.ParseFloat(v, 64)
}

// StrToInt конвертирует строку в int.
func StrToInt(v string) (int, error) {
	return strconv.Atoi(v)
}

// Compress сжимает данные при помощи пакета [compress/gzip].
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed init compress writer: %v", err)
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	return b.Bytes(), nil
}

// Decompress распаковывает сжатые данные при помощи пакета [compress/gzip].
func Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}

	return b.Bytes(), nil
}

// CalcSHA256 вычисляет хеш SHA-256 от переданных значения и ключа.
func CalcSHA256(value []byte, key string) (string, error) {
	h := sha256.New()
	_, err := h.Write(value)
	if err != nil {
		return "", err
	}
	return string(h.Sum(nil)), nil
}
