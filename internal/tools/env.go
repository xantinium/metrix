package tools

import (
	"os"
	"strings"
)

// IntEnvVar переменная окружения тип int.
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

// StrEnvVar переменная окружения тип string.
type StrEnvVar struct {
	Value  string
	Exists bool
}

// GetStrFromEnv достаёт переменную окружения типа string.
func GetStrFromEnv(name string) StrEnvVar {
	valueStr, exists := os.LookupEnv(name)
	if !exists {
		return StrEnvVar{}
	}

	return StrEnvVar{Exists: true, Value: valueStr}
}

// BoolEnvVar переменная окружения тип bool.
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
