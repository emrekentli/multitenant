package config

type Token struct {
	Hash         string
	Expire       int64
	AppJwtSecret string
	ApiJwtSecret string
}
