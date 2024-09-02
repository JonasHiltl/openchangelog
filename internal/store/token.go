package store

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/rs/xid"
)

const (
	token_prefix = "tkn"
)

type Token string

func NewToken() Token {
	id := xid.New()
	hasher := md5.New()
	hasher.Write(id.Bytes())

	return Token("tkn_" + hex.EncodeToString(hasher.Sum(nil)))
}

var errKeyFormat = errs.NewError(errs.ErrBadRequest, errors.New("wrong token key format"))

func ParseToken(key string) (Token, error) {
	parts := strings.Split(key, id_separator)
	if len(parts) != 2 {
		return "", errKeyFormat
	}
	if parts[0] != token_prefix {
		return "", errs.NewError(errs.ErrBadRequest, errors.New("invalid token key prefix"))
	}
	_, err := xid.FromString(parts[1])
	if err != nil {
		return "", errKeyFormat
	}
	return Token(key), nil
}

func (k Token) String() string {
	return string(k)
}

func (k Token) IsSet() bool {
	return string(k) != ""
}
