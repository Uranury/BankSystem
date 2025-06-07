package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtkey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(id int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtkey)
}

func VerifyJWT(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtkey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid id claim")
		}
		id := int64(userIDFloat)
		return id, nil
	}

	return 0, errors.New("invalid token")
}
