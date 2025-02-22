// Пакет tools содержит вспомогательные методы
// и структуры общего назначения.
package tools

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"strconv"
	"strings"
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

type IntEnvVar struct {
	Exists bool
	Value  int
}

// GetIntFromEnv достаёт переменную окружения типа int.
func GetIntFromEnv(name string) IntEnvVar {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		return IntEnvVar{}
	}

	value, err := StrToInt(valueStr)
	if err != nil {
		return IntEnvVar{}
	}

	return IntEnvVar{Exists: true, Value: value}
}

type StrEnvVar struct {
	Exists bool
	Value  string
}

// GetStrFromEnv достаёт переменную окружения типа string.
func GetStrFromEnv(name string) StrEnvVar {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		return StrEnvVar{}
	}

	return StrEnvVar{Exists: true, Value: valueStr}
}

type BoolEnvVar struct {
	Exists bool
	Value  bool
}

// GetBoolFromEnv достаёт переменную окружения типа bool.
func GetBoolFromEnv(name string) BoolEnvVar {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		return BoolEnvVar{}
	}

	return BoolEnvVar{Exists: true, Value: strings.TrimSpace(strings.ToLower(valueStr)) == "true"}
}
