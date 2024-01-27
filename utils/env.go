package utils

import "os"

func EnsureEnv(key, value string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		return value
	}
	return env
}
