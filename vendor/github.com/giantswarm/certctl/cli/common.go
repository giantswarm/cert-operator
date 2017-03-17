package cli

import (
	"os"
)

func fromEnv(key, def string) string {
	value := os.Getenv(key)

	if value == "" {
		return def
	}

	return value
}
