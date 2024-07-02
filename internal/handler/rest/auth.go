package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
)

type Token struct {
	Key         string
	WorkspaceID store.WorkspaceID
}

func bearerAuth(e *env, r *http.Request) (Token, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return Token{}, errs.NewError(errs.ErrUnauthorized, errors.New("missing authorization header"))
	}

	parts := strings.Split(h, " ")
	if len(parts) != 2 {
		return Token{}, errs.NewError(errs.ErrUnauthorized, errors.New("invalid bearer token format"))
	}
	key := parts[1]
	id, err := e.store.GetWorkspaceIDByToken(r.Context(), key)
	if err != nil {
		return Token{}, err
	}
	return Token{
		Key:         key,
		WorkspaceID: id,
	}, nil
}
