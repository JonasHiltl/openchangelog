package api

import (
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Token struct {
	Key         string
	WorkspaceID string
}

func (a *api) bearerAuth(c echo.Context) (Token, error) {
	h := c.Request().Header.Get("Authorization")
	if h == "" {
		return Token{}, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Missing Authorization header")
	}

	parts := strings.Split(h, " ")
	if len(parts) != 2 {
		return Token{}, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Invalid Bearer Token format")
	}
	key := parts[1]
	s, err := a.queries.GetToken(c.Request().Context(), key)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Token{}, echo.NewHTTPError(echo.ErrUnauthorized.Code, "Invalid Bearer Token")
		}
		return Token{}, err
	}
	return Token{
		Key:         s.Key,
		WorkspaceID: s.WorkspaceID,
	}, nil
}
