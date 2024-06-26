package apitypes

import (
	"encoding/json"
	"time"

	"github.com/guregu/null/v5"
)

// Represents the Changelog returned by the API via json encoding.
// Implements json un-/marshaling.
type Changelog struct {
	ID          int64     `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Title       string    `json:"title,omitempty"`
	Subtitle    string    `json:"subtitle,omitempty"`
	Logo        Logo      `json:"logo"`
	Source      Source    `json:"source,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (c *Changelog) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	// We'll store the error (if any) so we can return it if necessary
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	if idRaw, ok := objMap["id"]; ok {
		err = json.Unmarshal(*idRaw, &c.ID)
		if err != nil {
			return err
		}
	}

	if workspaceIdRaw, ok := objMap["workspaceId"]; ok {
		err = json.Unmarshal(*workspaceIdRaw, &c.WorkspaceID)
		if err != nil {
			return err
		}
	}

	if titleRaw, ok := objMap["title"]; ok {
		err = json.Unmarshal(*titleRaw, &c.Title)
		if err != nil {
			return err
		}
	}

	if subtitleRaw, ok := objMap["subtitle"]; ok {
		err = json.Unmarshal(*subtitleRaw, &c.Subtitle)
		if err != nil {
			return err
		}
	}

	if logoRaw, ok := objMap["logo"]; ok {
		err = json.Unmarshal(*logoRaw, &c.Logo)
		if err != nil {
			return err
		}
	}

	if createdAtRaw, ok := objMap["createdAt"]; ok {
		err = json.Unmarshal(*createdAtRaw, &c.CreatedAt)
		if err != nil {
			return err
		}
	}

	if sourceRaw, ok := objMap["source"]; ok && sourceRaw != nil {
		var sourceMap map[string]json.RawMessage
		err = json.Unmarshal(*sourceRaw, &sourceMap)
		if err != nil {
			return err
		}

		typeRaw, ok := sourceMap["type"]
		if !ok {
			// No source type specified, so no source is set.
			return nil
		}

		var Type string
		err = json.Unmarshal(typeRaw, &Type)
		if err != nil {
			return err
		}

		switch Type {
		case string(GitHub):
			var ghSource GHSource
			err = json.Unmarshal(*sourceRaw, &ghSource)
			if err != nil {
				return err
			}
			c.Source = ghSource
		}
	}

	return nil
}

type Logo struct {
	Src    null.String `json:"src"`
	Link   null.String `json:"link"`
	Alt    null.String `json:"alt"`
	Height null.String `json:"height"`
	Width  null.String `json:"width"`
}

func (l Logo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Src    string `json:"src,omitempty"`
		Link   string `json:"link,omitempty"`
		Alt    string `json:"alt,omitempty"`
		Height string `json:"height,omitempty"`
		Width  string `json:"width,omitempty"`
	}{
		Src:    l.Src.ValueOrZero(),
		Link:   l.Link.ValueOrZero(),
		Alt:    l.Alt.ValueOrZero(),
		Height: l.Height.ValueOrZero(),
		Width:  l.Width.ValueOrZero(),
	})
}

type CreateChangelogBody struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Logo     Logo   `json:"logo"`
}

type UpdateChangelogBody CreateChangelogBody
