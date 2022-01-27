package config

import (
	"os"
	"strings"
)

//Configuration is configuration struct
type Configuration struct {
	Port string
}

//GetConfig is configuration of the project
func GetConfig() Configuration {
	var c Configuration

	c.Port = getenv("PORT", ":8081")

	return c
}

func getenv(key, fallback string) string {
	// log.Println(os.Getenv(key))
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return strings.TrimSpace(value)
}
