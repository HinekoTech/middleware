package common

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func EncodedJwtToken(c *fiber.Ctx, res string) (string, error) {
	authorizeToken := c.Get("Authorization")
	if authorizeToken == "" {
		return "", errors.New("authorizeToken = nil ")

	}
	tokenString := strings.Replace(authorizeToken, "Bearer ", "", 1)
	if tokenString == "" {
		return "", errors.New("tokenString = nil ")

	}
	secretKeyToken := []byte(os.Getenv("secretKey"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKeyToken, nil
	})
	if err != nil {
		return err.Error(), err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimToken := claims[res].(string)
		return claimToken, nil
	}

	return "", err
}
