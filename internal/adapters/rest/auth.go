package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Token struct {
	Key         string
	WorkspaceID string
}

func bearerAuth(e *env, r *http.Request) (Token, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return Token{}, RestError{Code: http.StatusUnauthorized, Err: errors.New("missing authorization header")}
	}

	parts := strings.Split(h, " ")
	if len(parts) != 2 {
		return Token{}, RestError{Code: http.StatusUnauthorized, Err: errors.New("invalid bearer token format")}
	}
	key := parts[1]
	id, err := e.workspaceSrv.GetWorkspaceIDByToken(r.Context(), key)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Token{}, RestError{Code: http.StatusUnauthorized, Err: errors.New("invalid bearer token")}
		}
		return Token{}, err
	}
	return Token{
		Key:         key,
		WorkspaceID: id.String(),
	}, nil
}
