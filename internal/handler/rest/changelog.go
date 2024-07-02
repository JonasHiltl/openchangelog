package rest

import (
	"encoding/json"
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

	c, err := e.store.CreateChangelog(r.Context(), store.Changelog{
		WorkspaceID: t.WorkspaceID,
		ID:          store.NewCID(),
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		LogoSrc:     req.Logo.Src,
		LogoLink:    req.Logo.Link,
		LogoAlt:     req.Logo.Alt,
		LogoHeight:  req.Logo.Height,
		LogoWidth:   req.Logo.Width,
	})
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

	c, err := e.store.UpdateChangelog(r.Context(), t.WorkspaceID, cId, store.UpdateChangelogArgs{
		Title:      req.Title,
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

	var req apitypes.SetChangelogSourceBody
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	switch req.SourceType {
	case apitypes.GitHub:
		ghID, err := store.ParseGHID(req.SourceID)
		if err != nil {
			return err
		}
		err = e.store.SetChangelogGHSource(r.Context(), t.WorkspaceID, cId, ghID)
		if err != nil {
			return err
		}
	default:
		return errs.NewError(errs.ErrBadRequest, fmt.Errorf("invalid source type: %s", req.SourceType))
	}

	return nil
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
