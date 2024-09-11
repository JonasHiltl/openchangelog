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
				ColorScheme: Dark,
				CreatedAt:   now,
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
				"colorScheme": "dark",
				"createdAt": "%s"
			}`, nowStr),
		},
		{
			input: Changelog{
				ID:          "cl_xxxx",
				WorkspaceID: "ws_xxxx",
				Title:       NewString("Test Title"),
				ColorScheme: Automatic,
				CreatedAt:   now,
			},
			expect: fmt.Sprintf(`{
				"id": "cl_xxxx",
				"workspaceId": "ws_xxxx",
				"title": "Test Title",
				"colorScheme": "automatic",
				"createdAt": "%s"
			}`, nowStr),
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
		b, err := json.Marshal(table.input)
		if err != nil {
			t.Error(err)
		}

		assert.JSONEq(t, table.expect, string(b))
	}
}

func TestChangelogUnmarshal(t *testing.T) {
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
			ColorScheme: Dark,
			CreatedAt:   time.Unix(1715958564, 0).UTC(),
		},
		{
			ID:          "cl_xxxx",
			WorkspaceID: "ws_xxxx",
			Subdomain:   "workspace_fjhla",
			Title:       NewString("Test Title"),
			Subtitle:    NewString("Test Subtitle"),
			CreatedAt:   time.Unix(1715958564, 0).UTC(),
			ColorScheme: Light,
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

		assert.Equal(t, table, c)
	}
}

func TestUpdateChangelogBodyMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    UpdateChangelogBody
		expected string
	}{
		{
			name:  "empty struct",
			input: UpdateChangelogBody{},
			expected: `{
				"title": "",
				"subtitle": "",
				"logo": {},
				"domain": "",
				"subdomain": "",
				"colorScheme": ""
			}`,
		},
		{
			name: "valid title",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Title: NewString("test"),
				},
			},
			expected: `{
				"title": "test",
				"subtitle": "",
				"logo": {},
				"domain": "",
				"subdomain": "",
				"colorScheme": ""
			}`,
		},
		{
			name: "null title",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Title: NewNullString(),
				},
			},
			expected: `{
				"title": null,
				"subtitle": "",
				"logo": {},
				"domain": "",
				"subdomain": "",
				"colorScheme": ""
			}`,
		},
		{
			name: "valid logo",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Logo: Logo{
						Src: NewString("test"),
					},
				},
			},
			expected: `{
				"title": "",
				"subtitle": "",
				"logo": {
					"src": "test"
				},
				"domain": "",
				"subdomain": "",
				"colorScheme": ""
			}`,
		},
		{
			name: "valid color scheme",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					ColorScheme: Dark,
				},
			},
			expected: `{
				"title": "",
				"subtitle": "",
				"logo": {},
				"domain": "",
				"subdomain": "",
				"colorScheme": "dark"
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.input)
			if err != nil {
				t.Error(err)
			}

			assert.JSONEq(t, test.expected, string(b))
		})
	}
}

func TestUpdateChangelogBodyUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		input UpdateChangelogBody
	}{
		{
			name:  "empty body",
			input: UpdateChangelogBody{},
		},
		{
			name: "valid title",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Title: NewString("test"),
				},
			},
		},
		{
			name: "null title",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Title: NewNullString(),
				},
			},
		},
		{
			name: "valid logo",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Logo: Logo{
						Src: NewString("test"),
					},
				},
			},
		},
		{
			name: "null logo src",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					Logo: Logo{
						Src: NewNullString(),
					},
				},
			},
		},
		{
			name: "valid color scheme",
			input: UpdateChangelogBody{
				CreateChangelogBody: CreateChangelogBody{
					ColorScheme: Dark,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := json.Marshal(test.input)
			if err != nil {
				t.Error(err)
			}

			var body UpdateChangelogBody
			err = json.Unmarshal(b, &body)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, test.input, body)
		})
	}
}
