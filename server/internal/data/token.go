package data

import (
	"errors"
	"math/rand/v2"
)

const (
	TokenCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	TokenLength  = 9
)

var (
	ErrCannotGenerateUniqueToken = errors.New("cannot generate a unique token")
)

func GenerateRandomToken() string {
	token := make([]byte, TokenLength)
	for i := range token {
		token[i] = TokenCharSet[rand.IntN(len(TokenCharSet))]
	}
	return string(token)
}
