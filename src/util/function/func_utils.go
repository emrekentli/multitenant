package function

import (
	"log"
	"os"
	"strconv"
)

func GetEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Invalid value for %s: %v. Using default value %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

func StringToInt(connections string) int32 {
	connectionsInt, err := strconv.Atoi(connections)
	if err != nil {
		log.Println("Error converting string to int: ", err)
	}
	return int32(connectionsInt)
}
