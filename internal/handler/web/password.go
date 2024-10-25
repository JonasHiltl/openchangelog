package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jonashiltl/openchangelog/internal/changelog"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/handler/web/views"
)

func passwordSubmit(e *env, w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return views.PasswordProtectionError(err.Error()).Render(r.Context(), w)
	}

	pw := r.FormValue("password")
	if pw == "" {
		return views.PasswordProtectionError("missing password").Render(r.Context(), w)
	}

	u, err := url.Parse(r.Header.Get("HX-Current-URL"))
	if err != nil {
		return views.PasswordProtectionError(err.Error()).Render(r.Context(), w)
	}

	page, pageSize := handler.ParsePagination(u.Query())

	l, err := handler.LoadChangelog(e.loader, e.cfg.IsDBMode(), r, changelog.NewPagination(pageSize, page))
	if err != nil {
		return err
	}

	parsed := l.Parse(r.Context())

	err = handler.ValidatePassword(parsed.CL.PasswordHash, pw)
	if err != nil {
		return views.PasswordProtectionError(err.Error()).Render(r.Context(), w)
	}

	w.Header().Set("HX-Retarget", "body")
	// the hashed password does not add any actual security, but we do it for
	// obfuscation purposes
	setProtectedCookie(r, w, parsed.CL.PasswordHash)

	return e.render.RenderChangelog(r.Context(), w, RenderChangelogArgs{
		FeedURL:    handler.GetFeedURL(r),
		CurrentURL: handler.GetFullURL(r),
		CL:         parsed.CL,
		Articles:   parsed.Articles,
		HasMore:    parsed.HasMore,
		PageSize:   pageSize,
		NextPage:   page + 1,
	})
}

func setProtectedCookie(r *http.Request, w http.ResponseWriter, pwHash string) {
	const yearSeconds = 365 * 24 * 60 * 60

	c := &http.Cookie{
		Name:     createCookieKey(r),
		Value:    pwHash,
		MaxAge:   yearSeconds,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	// safari doesn't set secure cookie on localhost
	if getHost(r) == "localhost" {
		c.Secure = false
	}

	http.SetCookie(w, c)
}

func getProtectedCookieValue(r *http.Request) (string, error) {
	c, err := r.Cookie(createCookieKey(r))
	if err != nil {
		return "", err
	}

	return c.Value, nil
}

func getHost(r *http.Request) string {
	host := r.Host
	if r.Header.Get("X-Forwarded-Host") != "" {
		host = r.Header.Get("X-Forwarded-Host")
	}

	// remove port
	return strings.Split(host, ":")[0]
}

func createCookieKey(r *http.Request) string {
	wID, cID := handler.GetQueryIDs(r)
	if wID != "" && cID != "" {
		return fmt.Sprintf("protected-%s-%s", wID, cID)
	}

	host := getHost(r)

	return fmt.Sprintf("protected-%s", host)
}
