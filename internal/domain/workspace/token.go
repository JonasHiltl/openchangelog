package workspace

import (
	"errors"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/domain"
	"github.com/rs/xid"
)

type Token string

func NewToken() Token {
	return Token("tkn_" + xid.New().String())
}

func ParseToken(key string) (Token, error) {
	parts := strings.Split(key, "_")
	if len(parts) != 2 {
		return "", domain.NewError(domain.ErrBadRequest, errors.New("wrong token key format"))
	}
	if parts[0] != "tkn" {
		return "", domain.NewError(domain.ErrBadRequest, errors.New("invalid token key prefix"))
	}
	_, err := xid.FromString(parts[1])
	if err != nil {
		return "", domain.NewError(domain.ErrBadRequest, err)
	}
	return Token(key), nil
}

func (k Token) String() string {
	return string(k)
}

func (k Token) IsSet() bool {
	return string(k) != ""
}
