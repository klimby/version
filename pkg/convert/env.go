package convert

import "os"

// EnvToString Simple helper function to read an environment or return a default value.
func EnvToString(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
