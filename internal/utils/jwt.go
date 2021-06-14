package utils

import (
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	UserId      uint
	RandomBytes []byte // TODO: посмотреть как можно заменить
	jwt.StandardClaims
}

var (
	IncorrectToken = errors.New("Не верный токен")
	ForbidenToken  = errors.New("Токен больше не действителен")
)

// TODO: переделать пароль для подписи jwt
func CreateToken(id uint, bytes []byte) string {
	tk := &Token{UserId: id, RandomBytes: bytes}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("gfhjkm"))
	return tokenString
}

func VerifyToken(token string) (*Token, error) {
	tk := &Token{}

	tok, err := jwt.ParseWithClaims(token, tk, func(t *jwt.Token) (interface{}, error) {
		return []byte("gfhjkm"), nil
	})
	if err != nil {
		return nil, IncorrectToken
	}

	if !tok.Valid {
		return nil, ForbidenToken
	}

	return tk, nil
}
