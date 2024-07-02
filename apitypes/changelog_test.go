package apitypes

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/guregu/null/v5"
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
				Title:       null.NewString("Test Title", true),
				Subtitle:    null.NewString("Test Subtitle", true),
				Logo: Logo{
					Src:    null.NewString("logo src", true),
					Link:   null.NewString("logo link", true),
					Alt:    null.NewString("logo description", true),
					Height: null.NewString("30px", true),
					Width:  null.NewString("40px", true),
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
				"title": "Test Title",
				"subtitle": "Test Subtitle",
				"logo": {
					"src": "logo src",
					"link": "logo link",	
					"alt": "logo description",
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
	}

	for _, table := range tables {
		b, err := json.MarshalIndent(table.input, "", "\t")
		if err != nil {
			t.Error(err)
		}

		output := strings.Fields(string(b))
		expect := strings.Fields(table.expect)

		eq := reflect.DeepEqual(output, expect)
		if !eq {
			t.Errorf("Expected %s to equal %s", output, expect)
		}
	}
}

func TestChangelogUnmarshaling(t *testing.T) {
	tables := []Changelog{
		{
			ID:          "cl_xxxx",
			WorkspaceID: "ws_xxxx",
			Title:       null.NewString("Test Title", true),
			Subtitle:    null.NewString("Test Subtitle", true),
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
			Title:       null.NewString("Test Title", true),
			Subtitle:    null.NewString("Test Subtitle", true),
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

		eq := reflect.DeepEqual(table, c)
		if !eq {
			t.Errorf("Expected %+v to equal %+v", c, table)
		}
	}
}
