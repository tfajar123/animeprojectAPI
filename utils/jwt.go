package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "secret"

func GenerateToken(email string, userID int32) (string, error) {
	token := jwt.NewWithClaims((jwt.SigningMethodHS256), jwt.MapClaims{
		"email": email,
		"user_id": userID,
		"exp" : time.Now().Add(time.Hour * 168).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int32, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Cannot parse Token")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("Invalid Token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("Invalid Token Claims")
	}

	userID := int32(claims["user_id"].(float64))

	return int32(userID), nil
}