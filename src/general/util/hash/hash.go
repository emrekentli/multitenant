package hash

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func Hash(password string) string {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatalf("Hash oluşturulurken hata: %v", err)
	}
	return hash
}

func Match(password string, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatalf("Şifre doğrulama sırasında hata: %v", err)
	}
	return match
}
