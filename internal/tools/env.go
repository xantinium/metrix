package tools

import (
	"os"
	"strings"
)

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
