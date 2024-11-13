package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/handler"
	"github.com/jonashiltl/openchangelog/internal/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	changelog_id_param = "cid"
)

func changelogToApiType(cl store.Changelog) apitypes.Changelog {
	c := apitypes.Changelog{
		ID:          cl.ID.String(),
		Subdomain:   cl.Subdomain.String(),
		Domain:      cl.Domain.NullString(),
		WorkspaceID: cl.WorkspaceID.String(),
		Title:       cl.Title,
		Subtitle:    cl.Subtitle,
		Logo: apitypes.Logo{
			Src:    cl.LogoSrc,
			Link:   cl.LogoLink,
			Alt:    cl.LogoAlt,
			Height: cl.LogoHeight,
			Width:  cl.LogoWidth,
		},
		ColorScheme:   cl.ColorScheme.ToApiTypes(),
		HidePoweredBy: cl.HidePoweredBy,
		CreatedAt:     cl.CreatedAt,
		Protected:     cl.Protected,
		HasPassword:   cl.PasswordHash != "",
		Analytics:     cl.Analytics,
	}

	if cl.GHSource.Valid {
		c.Source = ghToApiType(cl.GHSource.ValueOrZero())
	}
	return c
}

func encodeChangelog(w http.ResponseWriter, cl store.Changelog) error {
	res := changelogToApiType(cl)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func createChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	var req apitypes.CreateChangelogBody
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	ws, err := e.store.GetWorkspace(r.Context(), t.WorkspaceID)
	if err != nil {
		return err
	}

	cl := store.Changelog{
		WorkspaceID:   t.WorkspaceID,
		ID:            store.NewCID(),
		Subdomain:     store.NewSubdomain(ws.Name),
		Title:         req.Title,
		Subtitle:      req.Subtitle,
		LogoSrc:       req.Logo.Src,
		LogoLink:      req.Logo.Link,
		LogoAlt:       req.Logo.Alt,
		LogoHeight:    req.Logo.Height,
		LogoWidth:     req.Logo.Width,
		HidePoweredBy: req.HidePoweredBy,
		Protected:     req.Protected,
		Analytics:     req.Analytics,
	}

	if req.ColorScheme == "" {
		cl.ColorScheme = store.System
	} else {
		cl.ColorScheme = store.NewColorScheme(req.ColorScheme)
	}

	if req.Password != "" {
		hash, err := hashPassword(req.Password)
		if err != nil {
			return errs.NewBadRequest(err)
		}
		cl.PasswordHash = hash
	}

	d, err := store.ParseDomainNullString(req.Domain)
	if err != nil {
		return err
	}
	cl.Domain = d

	c, err := e.store.CreateChangelog(r.Context(), cl)
	if err != nil {
		return err
	}
	return encodeChangelog(w, c)
}

func updateChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	var req apitypes.UpdateChangelogBody
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	cId, err := store.ParseCID(r.PathValue(changelog_id_param))
	if err != nil {
		return err
	}

	domain, err := store.ParseDomainNullString(req.Domain)
	if err != nil {
		return err
	}

	hashedPassword := req.Password
	if req.Password.IsValid() {
		// if password is actually defined, we hash it
		hash, err := hashPassword(req.Password.V())
		if err != nil {
			return errs.NewBadRequest(err)
		}
		hashedPassword = apitypes.NewString(hash)
	}

	c, err := e.store.UpdateChangelog(r.Context(), t.WorkspaceID, cId, store.UpdateChangelogArgs{
		Title:         req.Title,
		Subdomain:     req.Subdomain,
		Domain:        domain,
		Subtitle:      req.Subtitle,
		LogoSrc:       req.Logo.Src,
		LogoLink:      req.Logo.Link,
		LogoAlt:       req.Logo.Alt,
		LogoHeight:    req.Logo.Height,
		LogoWidth:     req.Logo.Width,
		ColorScheme:   store.NewColorScheme(req.ColorScheme),
		HidePoweredBy: req.HidePoweredBy,
		Protected:     req.Protected,
		PasswordHash:  hashedPassword,
		Analytics:     req.Analytics,
	})
	if err != nil {
		return err
	}
	return encodeChangelog(w, c)
}

func setChangelogSource(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := store.ParseCID(r.PathValue(changelog_id_param))
	if err != nil {
		return err
	}

	sId := r.PathValue("sid")
	if sId == "" {
		return errs.NewError(errs.ErrBadRequest, errors.New("missing sid path param"))
	}

	if store.IsGHID(sId) {
		ghID, err := store.ParseGHID(sId)
		if err != nil {
			return err
		}
		err = e.store.SetChangelogGHSource(r.Context(), t.WorkspaceID, cId, ghID)
		if err != nil {
			return err
		}
	} else {
		return errs.NewError(errs.ErrBadRequest, fmt.Errorf("invalid source id: %s", sId))
	}

	return nil
}

func deleteChangelogSource(e *env, _ http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := store.ParseCID(r.PathValue(changelog_id_param))
	if err != nil {
		return err
	}

	return e.store.DeleteChangelogSource(r.Context(), t.WorkspaceID, cId)
}

func getChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := store.ParseCID(r.PathValue(changelog_id_param))
	if err != nil {
		return err
	}

	c, err := e.store.GetChangelog(r.Context(), t.WorkspaceID, cId)
	if err != nil {
		return err
	}
	return encodeChangelog(w, c)
}

func getFullChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cID, err := store.ParseCID(r.PathValue(changelog_id_param))
	if err != nil {
		return errs.NewBadRequest(err)
	}

	cl, err := e.store.GetChangelog(r.Context(), t.WorkspaceID, cID)
	if err != nil {
		return errs.NewBadRequest(err)
	}

	page, pageSize := handler.ParsePagination(r.URL.Query())
	pagination := internal.NewPagination(pageSize, page)

	loaded, err := e.loader.LoadReleaseNotes(r.Context(), cl, pagination)
	if err != nil {
		return errs.NewBadRequest(err)
	}

	parsed := e.parser.Parse(r.Context(), loaded.Notes.Raw, pagination)

	articles := make([]apitypes.Article, len(parsed.Articles))
	for i, a := range parsed.Articles {
		content, _ := io.ReadAll(a.Content)

		articles[i] = apitypes.Article{
			ID:          a.Meta.ID,
			Title:       a.Meta.Title,
			Description: a.Meta.Description,
			PublishedAt: a.Meta.PublishedAt,
			Tags:        a.Meta.Tags,
			HTMLContent: string(content),
		}
	}
	res := apitypes.FullChangelog{
		Changelog:       changelogToApiType(loaded.CL),
		Articles:        articles,
		HasMoreArticles: parsed.HasMore,
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func listChangelogs(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cls, err := e.store.ListChangelogs(r.Context(), t.WorkspaceID)
	if err != nil {
		return err
	}

	res := make([]apitypes.Changelog, len(cls))
	for i, cl := range cls {
		res[i] = changelogToApiType(cl)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func deleteChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := store.ParseCID(r.PathValue(changelog_id_param))
	if err != nil {
		return err
	}

	err = e.store.DeleteChangelog(r.Context(), t.WorkspaceID, cId)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		return "", err
	}
	return string(hash), err
}
