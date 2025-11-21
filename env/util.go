package env

import (
	"os"
	"strconv"
)

// Loads value from env
// Panics if key does not exists in environ
func LoadEnvironStringOrPanic(key string) string {
	value, isExists := os.LookupEnv(key)

	if !isExists {
		panic("Failed to load environ variable by key: " + key)
	}

	return value
}

// Loads value from env and returns it, or defaultValue
func LoadEnvironStringWithDefault(key string, defaultValue string) string {
	value, isExists := os.LookupEnv(key)

	if !isExists {
		return defaultValue
	}

	return value
}

func ToIntOrPanic(value string) int {
	num, err := strconv.Atoi(value)

	if err != nil {
		panic("Failed to convert string to int: " + value)
	}

	return num
}
