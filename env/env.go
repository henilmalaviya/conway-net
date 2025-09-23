package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	Port                 string
	TickSpeed            int
	WSEndpoint           string
	WebSocketOriginCheck bool
	MaxObserveRegionSize int
	LogLevel             string
	SaveInterval         int
	SaveDirectory        string
	MaxSavesFiles        int
}

var env *Environment

func getEnvString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnvString(key, fmt.Sprintf("%d", fallback))
	intValue, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return intValue
}

func getEnvBool(key string, fallback bool) bool {
	valueStr := getEnvString(key, fmt.Sprintf("%t", fallback))
	boolValue, err := strconv.ParseBool(valueStr)
	if err != nil {
		return fallback
	}
	return boolValue
}

func Get() *Environment {
	return env
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	env = &Environment{
		Port:                 getEnvString("PORT", "8080"),
		TickSpeed:            getEnvInt("TICK_SPEED", 250),
		WSEndpoint:           getEnvString("WS_ENDPOINT", "/game"),
		WebSocketOriginCheck: getEnvBool("WS_ORIGIN_CHECK", false),
		MaxObserveRegionSize: getEnvInt("MAX_OBSERVE_REGION_SIZE", 1000),
		LogLevel:             getEnvString("LOG_LEVEL", "info"),
		SaveInterval:         getEnvInt("SAVE_INTERVAL", 60),
		SaveDirectory:        getEnvString("SAVE_DIR", "./saves"),
		MaxSavesFiles:        getEnvInt("MAX_SAVE_FILES", 10),
	}
}
