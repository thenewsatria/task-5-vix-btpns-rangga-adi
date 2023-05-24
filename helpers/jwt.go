package helpers

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IWebToken interface {
	GenerateAccessToken(email string) (string, error)
	ParseToken(tokenStr string) (*UserClaims, error)
}

type WebToken struct{}

type UserClaims struct {
	Email string
	jwt.RegisteredClaims
}

func NewWebToken() IWebToken {
	return &WebToken{}
}

func (wt *WebToken) GenerateAccessToken(email string) (string, error) {
	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		return "", nil
	}
	tokenSecret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(time.Duration(expTime) * time.Minute)
	claims := &UserClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (wt *WebToken) ParseToken(tokenStr string) (*UserClaims, error) {
	tokenSecret := os.Getenv("JWT_SECRET")
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("token: Token invalid")
	}
}
