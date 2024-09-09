package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jonashiltl/openchangelog/apitypes"
	"github.com/jonashiltl/openchangelog/internal/errs"
	"github.com/jonashiltl/openchangelog/internal/store"
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
		CreatedAt: cl.CreatedAt,
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
		WorkspaceID: t.WorkspaceID,
		ID:          store.NewCID(),
		Subdomain:   store.NewSubdomain(ws.Name),
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		LogoSrc:     req.Logo.Src,
		LogoLink:    req.Logo.Link,
		LogoAlt:     req.Logo.Alt,
		LogoHeight:  req.Logo.Height,
		LogoWidth:   req.Logo.Width,
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

	c, err := e.store.UpdateChangelog(r.Context(), t.WorkspaceID, cId, store.UpdateChangelogArgs{
		Title:      req.Title,
		Subdomain:  req.Subdomain,
		Domain:     domain,
		Subtitle:   req.Subtitle,
		LogoSrc:    req.Logo.Src,
		LogoLink:   req.Logo.Link,
		LogoAlt:    req.Logo.Alt,
		LogoHeight: req.Logo.Height,
		LogoWidth:  req.Logo.Width,
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
