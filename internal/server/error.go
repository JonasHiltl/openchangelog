package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if message, ok := he.Message.(string); ok {
			msg = message
		}
	}

	req := c.Request()
	path := req.URL.Path
	if req.URL.RawQuery != "" {
		path += "?" + req.URL.RawQuery
	}

	if strings.HasPrefix(path, "/api") {
		c.JSON(code, echo.HTTPError{
			Code:    code,
			Message: msg,
		})
	} else {
		vars := map[string]interface{}{
			"Status":  code,
			"Message": msg,
			"Path":    path,
		}
		if err := c.Render(code, "error", vars); err != nil {
			c.Logger().Error(err)
		}
	}
}
