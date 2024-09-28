package utils

import "os"

func GetEnvOrDefault(name, defaultValue string) string {
	if variable := os.Getenv(name); variable != "" {
		return variable
	}
	return defaultValue
}
