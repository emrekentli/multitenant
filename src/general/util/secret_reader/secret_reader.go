package secret_reader

import (
	"log"
	"os"
)

func ReadSecret(secretDir string) string {
	log.Println("Reading secret file: ", secretDir)
	secret, err := os.ReadFile(secretDir)
	if err != nil {
		log.Println("Error reading secret file: ", err)
		return secretDir
	}
	return string(secret)
}
