package config

import (
	"fmt"
	"os"
)

func MustGetEnv(key string) string {
	env := os.Getenv(key)
	if env == "" {
		panic(fmt.Sprintf("env %s not set", key))
	}
	return env
}
