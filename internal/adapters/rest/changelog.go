package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jonashiltl/openchangelog/internal/domain/changelog"
	"github.com/jonashiltl/openchangelog/internal/domain/source"
)

type Logo struct {
	Src    string `json:"src"`
	Link   string `json:"link"`
	Alt    string `json:"alt"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

func changelogToMap(cl changelog.Changelog) map[string]any {
	m := map[string]any{
		"id":          cl.ID,
		"workspaceID": cl.WorkspaceID,
		"createdAt":   cl.CreatedAt,
	}
	if cl.Title != "" {
		m["title"] = cl.Title
	}
	if cl.Subtitle != "" {
		m["subtitle"] = cl.Subtitle
	}

	if cl.Source != nil {
		switch cl.Source.Type() {
		case source.GitHub:
			g := cl.Source.(source.GHSource)
			m["source"] = ghToMap(g)
		}
	}

	if cl.Logo.Src != "" {
		m["logo"] = Logo(cl.Logo)
	}

	return m
}

func encodeChangelog(w http.ResponseWriter, cl changelog.Changelog) error {
	res := changelogToMap(cl)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func createChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}
	var req struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Logo     Logo   `json:"logo"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	c, err := e.changelogSrv.CreateChangelog(r.Context(), changelog.CreateChangelogArgs{
		WorkspaceID: t.WorkspaceID,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Logo: struct {
			Src    string
			Link   string
			Alt    string
			Height string
			Width  string
		}(req.Logo),
	})
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return encodeChangelog(w, c)
}

func updateChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}
	var req struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Logo     Logo   `json:"logo"`
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(r.PathValue("cid"), 10, 64)
	if err != nil {
		return err
	}

	c, err := e.changelogSrv.UpdateChangelog(
		r.Context(),
		t.WorkspaceID,
		cId,
		changelog.UpdateChangelogArgs{
			Title:    req.Title,
			Subtitle: req.Subtitle,
			Logo: struct {
				Src    string
				Link   string
				Alt    string
				Height string
				Width  string
			}(req.Logo),
		},
	)
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return encodeChangelog(w, c)
}

func setChangelogSource(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(r.PathValue("cid"), 10, 64)
	if err != nil {
		return err
	}

	var req struct {
		SourceType string `json:"type"`
		SourceID   int64  `json:"id"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = e.changelogSrv.SetChangelogSource(r.Context(), t.WorkspaceID, cId, changelog.SetChangelogSourceArgs{
		Type: req.SourceType,
		ID:   req.SourceID,
	})
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return nil
}

func getChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(r.PathValue("cid"), 10, 64)
	if err != nil {
		return err
	}

	c, err := e.changelogSrv.GetChangelog(r.Context(), t.WorkspaceID, cId)
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return encodeChangelog(w, c)
}

func listChangelogs(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cls, err := e.changelogSrv.ListChangelogs(r.Context(), t.WorkspaceID)
	if err != nil {
		return RestErrorFromDomain(err)
	}

	res := make([]map[string]any, len(cls))
	for i, cl := range cls {
		res[i] = changelogToMap(cl)
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func deleteChangelog(e *env, w http.ResponseWriter, r *http.Request) error {
	t, err := bearerAuth(e, r)
	if err != nil {
		return err
	}

	cId, err := strconv.ParseInt(r.PathValue("cid"), 10, 64)
	if err != nil {
		return err
	}

	err = e.changelogSrv.DeleteChangelog(r.Context(), t.WorkspaceID, cId)
	if err != nil {
		return RestErrorFromDomain(err)
	}
	return nil
}
