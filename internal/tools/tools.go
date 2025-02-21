// Пакет tools содержит вспомогательные методы
// и структуры общего назначения.
package tools

import (
	"bytes"
	"compress/flate"
	"fmt"
	"strconv"
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

// Compress сжимает данные при помощи пакета [compress/flate].
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := flate.NewWriter(&b, flate.BestCompression)
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
