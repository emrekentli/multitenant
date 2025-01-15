package jwt

import (
	"app/config"
	"app/src/general/util/rest"
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

var specialAuthorizedRoutes = []string{
	"/login",
}

func RegisterJwtMiddleware(app *fiber.App) {
	app.Use(func(c fiber.Ctx) error {
		if filterJwtMiddleware(c) {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			authHeader = c.Cookies("Authorization")
		}

		if authHeader == "" {
			return onUnAuthorized(c)
		}

		token, err := VerifyToken(authHeader)
		if err != nil || !token.Valid {
			return onUnAuthorized(c)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return onUnAuthorized(c)
		}
		c.Locals("ID", claims["ID"])

		return c.Next()
	})
}

func filterJwtMiddleware(c fiber.Ctx) bool {
	for _, route := range specialAuthorizedRoutes {
		if strings.Contains(c.Path(), route) {
			return true
		}
	}
	return false
}

func onUnAuthorized(c fiber.Ctx) error {
	return rest.ErrorRes(c, rest.Unauthorized, rest.ErrorCode[rest.Unauthorized])
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Token'ı parse et
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	// Token'ın expiration süresini kontrol et
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("token expired")
			}
		}
	}
	return token, nil
}

func CreateJwt(id int64) (string, error) {
	day := 24 * time.Hour
	claims := jwt.MapClaims{
		"ID":  id,
		"exp": time.Now().Add(day * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(config.JwtSecretByte)
	if err != nil {
		return "", err
	}
	return t, nil
}
