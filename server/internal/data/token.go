package data

import (
	"errors"
	"math/rand/v2"
	"regexp"

	"github.com/soumikc1729/splitty/server/internal/validator"
)

const (
	TokenCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	TokenLength  = 9
)

var (
	TokenFormatRX                = regexp.MustCompile("^[A-Z0-9]{9}$")
	ErrCannotGenerateUniqueToken = errors.New("cannot generate a unique token")
)

func ValidateToken(v *validator.Validator, token string) {
	v.Check(validator.Matches(token, TokenFormatRX), "token", "must be 9 characters long and contain only letters and numbers")
}

func GenerateRandomToken() string {
	token := make([]byte, TokenLength)
	for i := range token {
		token[i] = TokenCharSet[rand.IntN(len(TokenCharSet))]
	}
	return string(token)
}
