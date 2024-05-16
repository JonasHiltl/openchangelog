package apitypes

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
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
				ID:          1,
				WorkspaceID: "ws_xxxx",
				Title:       "Test Title",
				Subtitle:    "Test Subtitle",
				Source: GHSource{
					ID:          1,
					WorkspaceID: "ws_xxxx",
					Owner:       "jonashiltl",
					Repo:        "openchangelog",
					Path:        ".testdata",
				},
				CreatedAt: now,
			},
			expect: fmt.Sprintf(`{
				"id": 1,
				"workspaceId": "ws_xxxx",
				"title": "Test Title",
				"subtitle": "Test Subtitle",
				"logo": {
					"src": "",
					"link": "",
					"alt": "",
					"height": "",
					"width": ""
				},
				"source": {
					"type": "github",
					"id": 1,
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
