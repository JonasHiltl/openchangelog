package apitypes

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChangelogMarshaling(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339Nano)
	tables := []struct {
		input  Changelog
		expect string
	}{
		{
			input: Changelog{
				ID:          "cl_xxxx",
				WorkspaceID: "ws_xxxx",
				Subdomain:   "workspace_fjhla",
				Domain:      NewString("demo.openchangelog.com"),
				Title:       NewString("Test Title"),
				Subtitle:    NewString("Test Subtitle"),
				Logo: Logo{
					Src:    NewString("logo src"),
					Link:   NewString("logo link"),
					Alt:    NewString("logo description"),
					Height: NewString("30px"),
					Width:  NewString("40px"),
				},
				Source: GHSource{
					ID:          "gh_xxxx",
					WorkspaceID: "ws_xxxx",
					Owner:       "jonashiltl",
					Repo:        "openchangelog",
					Path:        ".testdata",
				},
				CreatedAt: now,
			},
			expect: fmt.Sprintf(`{
				"id": "cl_xxxx",
				"workspaceId": "ws_xxxx",
				"subdomain": "workspace_fjhla",
				"title": "Test Title",
				"domain": "demo.openchangelog.com",
				"subtitle": "Test Subtitle",
				"logo": {
					"alt": "logo description",
					"src": "logo src",
					"link": "logo link",	
					"height": "30px",
					"width": "40px"
				},
				"source": {
					"type": "github",
					"id": "gh_xxxx",
					"workspaceId": "ws_xxxx",
					"owner": "jonashiltl",
					"repo": "openchangelog",
					"path": ".testdata"
				},
				"createdAt": "%s"
			}`, nowStr),
		},
		{
			input: Changelog{
				ID:          "cl_xxxx",
				WorkspaceID: "ws_xxxx",
				Title:       NewString("Test Title"),
			},
			expect: `{
				"id": "cl_xxxx",
				"workspaceId": "ws_xxxx",
				"title": "Test Title"
			}`,
		},
		{
			input: Changelog{
				ID:          "cl_xxxx",
				WorkspaceID: "ws_xxxx",
				Logo: Logo{
					Alt: NewString("test"),
				},
			},
			expect: `{
				"id": "cl_xxxx",
				"workspaceId": "ws_xxxx",
				"logo": {
					"alt": "test"
				}
			}`,
		},
	}

	for _, table := range tables {
		b, err := json.MarshalIndent(table.input, "", "\t")
		if err != nil {
			t.Error(err)
		}

		assert.JSONEq(t, table.expect, string(b))
	}
}

func TestChangelogUnmarshaling(t *testing.T) {
	tables := []Changelog{
		{
			ID:          "cl_xxxx",
			WorkspaceID: "ws_xxxx",
			Subdomain:   "workspace_fjhla",
			Domain:      NewString("demo.openchangelog.com"),
			Title:       NewString("Test Title"),
			Subtitle:    NewString("Test Subtitle"),
			Source: GHSource{
				ID:          "gh_xxxx",
				WorkspaceID: "ws_xxxx",
				Owner:       "jonashiltl",
				Repo:        "openchangelog",
				Path:        ".testdata",
			},
			CreatedAt: time.Unix(1715958564, 0),
		},
		{
			ID:          "cl_xxxx",
			WorkspaceID: "ws_xxxx",
			Subdomain:   "workspace_fjhla",
			Title:       NewString("Test Title"),
			Subtitle:    NewString("Test Subtitle"),
			CreatedAt:   time.Unix(1715958564, 0),
		},
	}

	for _, table := range tables {
		b, err := json.Marshal(table)
		if err != nil {
			t.Error(err)
		}

		var c Changelog
		err = json.Unmarshal(b, &c)
		if err != nil {
			t.Error(err)
		}

		// expected logo isNull: false
		// actual logo isNull: true

		// fails because Logo is empty, not equal after unmarshaling
		assert.Equal(t, table, c)
	}
}
