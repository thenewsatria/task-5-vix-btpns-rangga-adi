package helpers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IWebToken interface {
	GenerateToken(userId uint) (string, error)
	ParseToken(tokenStr string) (*UserClaims, error)
	IsTokenExpired(tokenStr string) (bool, error)
}

type WebToken struct {
	expirationTimeInMinute int
	tokenSecret            string
}

type UserClaims struct {
	ID uint
	jwt.RegisteredClaims
}

func NewWebToken(expirationTimeInMinute int, tokenSecret string) IWebToken {
	return &WebToken{
		expirationTimeInMinute: expirationTimeInMinute,
		tokenSecret:            tokenSecret,
	}
}

func (wt *WebToken) GenerateToken(userId uint) (string, error) {
	expirationTime := time.Now().Add(time.Duration(wt.expirationTimeInMinute) * time.Minute)
	claims := &UserClaims{
		ID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(wt.tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (wt *WebToken) ParseToken(tokenStr string) (*UserClaims, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(wt.tokenSecret), nil
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

func (wt *WebToken) IsTokenExpired(tokenStr string) (bool, error) {
	claims := &UserClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(wt.tokenSecret), nil
	})
	if errors.Is(err, jwt.ErrTokenExpired) {
		return true, nil
	} else {
		return false, err
	}
}
