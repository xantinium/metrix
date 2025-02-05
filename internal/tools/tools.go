// Пакет tools содержит вспомогательные методы
// и структуры общего назначения.
package tools

import "strconv"

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
